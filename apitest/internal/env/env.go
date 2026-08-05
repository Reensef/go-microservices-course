// Package env задаёт адреса уже поднятого dev-стека (docker-compose), по
// которым бьют сценарии. Сценарии сами ничего не поднимают и не
// останавливают — предполагается, что стек уже запущен (task compose:*:up).
package env

import "os"

type Config struct {
	IAMAddr       string // gRPC-адрес iam (AuthService, UserService)
	InventoryAddr string // gRPC-адрес inventory
	PaymentAddr   string // gRPC-адрес payment
	OrderBaseURL  string // HTTP base URL order
}

func Load() Config {
	return Config{
		IAMAddr:       getEnv("APITEST_IAM_ADDR", "localhost:50054"),
		InventoryAddr: getEnv("APITEST_INVENTORY_ADDR", "localhost:50051"),
		PaymentAddr:   getEnv("APITEST_PAYMENT_ADDR", "localhost:50053"),
		OrderBaseURL:  getEnv("APITEST_ORDER_BASE_URL", "http://localhost:50052"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
