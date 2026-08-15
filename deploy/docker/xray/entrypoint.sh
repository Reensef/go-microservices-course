#!/bin/sh
# Точка входа xray-прокси.
#
# На каждом старте контейнера подписка скачивается заново, поэтому любой
# перезапуск (в том числе автоматический, по restart: unless-stopped)
# поднимает прокси с актуальным списком серверов.
#
# После запуска следим за связью: если несколько проверок подряд не прошли,
# процесс завершается, Docker перезапускает контейнер — и цикл начинается
# со свежей подписки.

set -eu

SUB_URL="${XRAY_SUB_URL:?XRAY_SUB_URL is not set}"
CONFIG_INDEX="${XRAY_CONFIG_INDEX:-1}"
HTTP_PORT="${XRAY_HTTP_PORT:-10809}"
PROBE_URL="${XRAY_PROBE_URL:-https://api.telegram.org}"
PROBE_INTERVAL="${XRAY_PROBE_INTERVAL:-30}"
PROBE_RETRIES="${XRAY_PROBE_RETRIES:-3}"

# Обрыв обнаруживается за PROBE_INTERVAL * (PROBE_RETRIES + 1): сон идёт и перед
# первой проверкой тоже. При 5 и 3 это 20 секунд, а не 15
CONFIG_PATH=/etc/xray/config.json
SUB_PATH=/tmp/subscription.json

# Потолок задержки между попытками скачать подписку. Определяет, сколько мы
# впустую прождём уже после того, как связь вернулась
MAX_RETRY_DELAY=8

log() {
	echo "[xray-entrypoint] $*"
}

# Скачиваем подписку с повторами: сеть после старта хоста может подняться не сразу.
#
# --connect-timeout нужен на случай, когда пакеты дропаются молча, без RST:
# без него попытка висела бы до --max-time, и это, а не задержки между
# попытками, было бы главным слагаемым времени восстановления
fetch_subscription() {
	attempt=1
	delay=2

	while [ "$attempt" -le 5 ]; do
		if curl -fsSL --connect-timeout 5 --max-time 20 "$SUB_URL" -o "$SUB_PATH"; then
			return 0
		fi

		log "failed to download subscription (attempt $attempt of 5), retrying in ${delay}s"
		sleep "$delay"
		attempt=$((attempt + 1))

		delay=$((delay * 2))
		if [ "$delay" -gt "$MAX_RETRY_DELAY" ]; then
			delay="$MAX_RETRY_DELAY"
		fi
	done

	return 1
}

# Берём нужный конфиг из подписки и заменяем inbounds единственным HTTP-прокси.
#
# listen 0.0.0.0 здесь безопасен: это 0.0.0.0 внутри сетевого namespace
# контейнера, а порт не публикуется через ports, поэтому прокси доступен
# только сервисам в сети common и не виден ни с хоста, ни снаружи.
build_config() {
	if ! jq -e --argjson idx "$CONFIG_INDEX" '.[$idx] | objects' "$SUB_PATH" >/dev/null; then
		log "subscription has no config at index $CONFIG_INDEX"
		return 1
	fi

	jq --argjson idx "$CONFIG_INDEX" --argjson port "$HTTP_PORT" '
		.[$idx]
		| .inbounds = [{
			tag: "http",
			listen: "0.0.0.0",
			port: $port,
			protocol: "http",
			settings: { allowTransparent: false }
		}]
		| del(.remarks)
	' "$SUB_PATH" >"$CONFIG_PATH"
}

log "downloading subscription"
if ! fetch_subscription; then
	log "subscription unavailable, exiting — docker will restart the container"
	exit 1
fi

if ! build_config; then
	exit 1
fi

log "config built, http proxy listening on 0.0.0.0:${HTTP_PORT}"

xray run -config "$CONFIG_PATH" &
xray_pid=$!

trap 'log "signal received, stopping xray"; kill -TERM "$xray_pid" 2>/dev/null || true; wait "$xray_pid" 2>/dev/null || true; exit 0' TERM INT

failures=0

while true; do
	# sleep в фоне + wait, иначе шелл не отреагирует на SIGTERM до конца паузы
	sleep "$PROBE_INTERVAL" &
	wait $! || true

	if ! kill -0 "$xray_pid" 2>/dev/null; then
		log "xray process exited, quitting for restart"
		exit 1
	fi

	if curl -fsS --max-time 15 -o /dev/null -x "http://127.0.0.1:${HTTP_PORT}" "$PROBE_URL"; then
		failures=0
		continue
	fi

	failures=$((failures + 1))
	log "connectivity probe failed ($failures of $PROBE_RETRIES)"

	if [ "$failures" -ge "$PROBE_RETRIES" ]; then
		log "connectivity lost, exiting — docker will restart the container with a fresh subscription"
		kill -TERM "$xray_pid" 2>/dev/null || true
		wait "$xray_pid" 2>/dev/null || true
		exit 1
	fi
done
