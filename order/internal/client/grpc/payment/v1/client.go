package v1

import (
	"sync"

	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

var _ grpcClients.PaymentServiceClient = (*paymentClient)(nil)

type paymentClient struct {
	mu      sync.RWMutex
	service paymentV1.PaymentServiceClient
}

func NewClient(service paymentV1.PaymentServiceClient) *paymentClient {
	return &paymentClient{
		service: service,
	}
}
