// Package clients собирает типизированные клиенты ко всем сервисам, которые
// дёргают сценарии, плюс общую для gRPC и HTTP авторизацию через контекст
// (см. WithSession в auth.go).
package clients

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Reensef/go-microservices-course/apitest/internal/env"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	paymentGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

// Clients — один набор соединений на весь тестовый процесс: gRPC-соединения
// и HTTP-клиент order дёшево переиспользовать между сценариями/шагами.
type Clients struct {
	IAMAuth   iamGrpc.AuthServiceClient
	IAMUser   iamGrpc.UserServiceClient
	Inventory inventoryGrpc.InventoryServiceClient
	Payment   paymentGrpc.PaymentServiceClient
	Order     *orderApi.Client
}

// New поднимает клиентские соединения к уже запущенному dev-стеку.
// Если какой-то адрес недоступен, тест падает сразу же с понятной причиной —
// как правило это значит, что стек не поднят (см. Taskfile: compose:*:up).
func New(t *testing.T) *Clients {
	t.Helper()

	cfg := env.Load()

	iamConn := dial(t, "iam", cfg.IAMAddr)
	inventoryConn := dial(t, "inventory", cfg.InventoryAddr)
	paymentConn := dial(t, "payment", cfg.PaymentAddr)

	orderClient, err := orderApi.NewClient(cfg.OrderBaseURL, orderApi.WithClient(&http.Client{
		Transport: bearerRoundTripper{base: http.DefaultTransport},
	}))
	require.NoError(t, err, "failed to build order HTTP client")

	return &Clients{
		IAMAuth:   iamGrpc.NewAuthServiceClient(iamConn),
		IAMUser:   iamGrpc.NewUserServiceClient(iamConn),
		Inventory: inventoryGrpc.NewInventoryServiceClient(inventoryConn),
		Payment:   paymentGrpc.NewPaymentServiceClient(paymentConn),
		Order:     orderClient,
	}
}

func dial(t *testing.T, name, addr string) *grpc.ClientConn {
	t.Helper()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoErrorf(t, err, "failed to dial %s at %s — is the dev stack up? (task compose:%s:up)", name, addr, name)

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close %s gRPC connection: %s", name, err)
		}
	})

	return conn
}
