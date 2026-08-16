# ============================
# Stage 1: Xray binary
# ============================

# Официальный образ xray-core — distroless, поэтому берём из него только
# бинарник и geo-базы, а запускать будем поверх alpine с curl и jq
FROM ghcr.io/xtls/xray-core:latest AS xray


# ============================
# Stage 2: Final image
# ============================

FROM alpine:3.21.3

# curl и jq нужны точке входа: скачать подписку, собрать из неё конфиг
# и проверять связь; ca-certificates — для HTTPS
RUN apk add --no-cache ca-certificates curl jq

# Бинарник и geo-базы (geoip.dat / geosite.dat) — их требуют routing-правила подписки
COPY --from=xray /usr/local/bin/xray /usr/local/bin/xray
COPY --from=xray /usr/local/share/xray/ /usr/local/share/xray/

COPY deploy/docker/xray/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Создаём системного пользователя без root-прав — best practice для безопасности.
# Конфиг собирается на старте, поэтому каталог должен быть ему доступен на запись
RUN addgroup -S xray && adduser -S xray -G xray \
	&& mkdir -p /etc/xray \
	&& chown xray:xray /etc/xray

USER xray

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
