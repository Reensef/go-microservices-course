package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) PayOrder(
	ctx context.Context,
	orderUuid uuid.UUID,
	userUuid uuid.UUID,
	paymentMethod model.OrderPaymentMethod,
) (*uuid.UUID, error) {
	transactionUuid, err := s.paymentService.PayOrder(ctx, orderUuid, userUuid, paymentMethod)
	if err != nil {
		return nil, err
	}

	err = s.orderRepo.PayOrder(ctx, orderUuid, *transactionUuid, paymentMethod)
	if err != nil {
		return nil, err
	}

	return transactionUuid, nil
}
