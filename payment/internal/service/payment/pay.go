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
	if paymentMethod == model.PaymentMethod_UNSPECIFIED {
		return nil, model.ErrPaymentMethodUnspecified
	}

	transactionUuid := uuid.NewString()

	log.Println("Оплата прошла успешно, transaction_uuid:", transactionUuid)

	return &transactionUuid, nil
}
