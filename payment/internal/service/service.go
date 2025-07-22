package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

type PaymentService interface {
	Pay(
		ctx context.Context,
		orderUuid, userUuid uuid.UUID,
		paymentMethod model.PaymentMethod,
	) (*uuid.UUID, error)
}
