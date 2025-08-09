package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

func (s *service) Pay(
	ctx context.Context,
	orderUuid, userUuid uuid.UUID,
	paymentMethod model.PaymentMethod,
) (*uuid.UUID, error) {
	if paymentMethod == model.PaymentMethod_UNSPECIFIED {
		return nil, model.ErrPaymentMethodUnspecified
	}

	transactionUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	log.Println("Оплата прошла успешно, transaction_uuid:", transactionUuid.String())

	return &transactionUuid, nil
}
