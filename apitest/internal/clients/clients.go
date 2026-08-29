// Package clients собирает типизированные клиенты ко всем сервисам, которые
// дёргают сценарии — все они HTTP-клиенты одного Envoy-гейтвея, плюс общая
// авторизация через контекст (см. WithSession в auth.go).
package clients

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Reensef/go-microservices-course/apitest/internal/env"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

// Clients — один набор HTTP-клиентов на весь тестовый процесс.
type Clients struct {
	IAM       *iamClient
	Inventory *inventoryClient
	Order     *orderApi.Client
}

// New поднимает клиенты к уже запущенному dev-стеку (через Envoy-гейтвей).
func New(t *testing.T) *Clients {
	t.Helper()

	cfg := env.Load()

	httpClient := &http.Client{
		Transport: bearerRoundTripper{base: http.DefaultTransport},
	}

	orderClient, err := orderApi.NewClient(cfg.GatewayBaseURL, orderApi.WithClient(httpClient))
	require.NoError(t, err, "failed to build order HTTP client")

	return &Clients{
		IAM:       &iamClient{baseURL: cfg.GatewayBaseURL, http: httpClient},
		Inventory: &inventoryClient{baseURL: cfg.GatewayBaseURL, http: httpClient},
		Order:     orderClient,
	}
}
