package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

type PaymentService interface {
	Pay(
		ctx context.Context,
		orderUuid, userUuid string,
		paymentMethod model.PaymentMethod,
	) (*string, error)
}
