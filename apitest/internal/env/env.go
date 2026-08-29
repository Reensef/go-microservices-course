// Package env задаёт адрес уже поднятого dev-стека (docker-compose), по
// которому бьют сценарии. Сценарии сами ничего не поднимают и не
// останавливают — предполагается, что стек уже запущен (task compose:*:up
// плюс task compose:envoy:up).
package env

import "os"

type Config struct {
	GatewayBaseURL string // HTTP base URL Envoy-гейтвея — единственная точка входа
}

func Load() Config {
	return Config{
		GatewayBaseURL: getEnv("APITEST_GATEWAY_BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
