package v1

import (
	"sync"

	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
)

var _ grpcClients.PaymentServiceClient = (*paymentClient)(nil)

type paymentClient struct {
	mu sync.RWMutex
}

func NewClient() *paymentClient {
	return &paymentClient{}
}
