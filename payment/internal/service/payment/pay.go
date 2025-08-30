package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

func (s *service) Pay(
	ctx context.Context,
	orderUuid, userUuid string,
	paymentMethod model.PaymentMethod,
) (*string, error) {
	if uuid.Validate(orderUuid) != nil {
		return nil, model.ErrOrderUuidInvalidFormat
	}
	if uuid.Validate(userUuid) != nil {
		return nil, model.ErrUserUuidInvalidFormat
	}
	if paymentMethod == model.PaymentMethod_UNSPECIFIED {
		return nil, model.ErrPaymentMethodUnspecified
	}

	transactionUuid := uuid.NewString()

	return &transactionUuid, nil
}
