package app

import (
	"context"

	paymentApi "github.com/Reensef/go-microservices-course/payment/internal/api/payment/v1"
	service "github.com/Reensef/go-microservices-course/payment/internal/service"
	paymentService "github.com/Reensef/go-microservices-course/payment/internal/service/payment"
	paymentProtoApi "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentApi     paymentProtoApi.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentApi(ctx context.Context) paymentProtoApi.PaymentServiceServer {
	if d.paymentApi == nil {
		d.paymentApi = paymentApi.New(d.PaymentService(ctx))
	}

	return d.paymentApi
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.New()
	}

	return d.paymentService
}
