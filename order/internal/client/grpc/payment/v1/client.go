package v1

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	paymentGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

var _ grpcClients.PaymentClient = (*paymentClient)(nil)

type paymentClient struct {
	service paymentGrpc.PaymentServiceClient
}

func New(service paymentGrpc.PaymentServiceClient) *paymentClient {
	return &paymentClient{
		service: service,
	}
}
