package v1

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

var _ grpcClients.PaymentServiceClient = (*paymentV1Client)(nil)

type paymentV1Client struct {
	service paymentV1.PaymentServiceClient
}

func NewClient(service paymentV1.PaymentServiceClient) *paymentV1Client {
	return &paymentV1Client{
		service: service,
	}
}
