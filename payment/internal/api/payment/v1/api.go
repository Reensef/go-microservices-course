package v1

import (
	service "github.com/Reensef/go-microservices-course/payment/internal/service"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	service service.PaymentService
}

func NewAPI(service service.PaymentService) *api {
	return &api{
		service: service,
	}
}
