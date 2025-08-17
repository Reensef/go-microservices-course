package order

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) PayOrder(
	ctx context.Context,
	orderUuid string,
	userUuid string,
	paymentMethod model.OrderPaymentMethod,
) (*string, error) {
	status, err := s.orderRepo.GetOrderStatus(ctx, orderUuid)
	if err != nil {
		return nil, err
	}

	if status == model.OrderStatus_PAID {
		return nil, model.ErrOrderAlreadyPaid
	}

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
