package payment

import (
	"context"
	"log"

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

	log.Println("Оплата прошла успешно, transaction_uuid:", transactionUuid)

	return &transactionUuid, nil
}
