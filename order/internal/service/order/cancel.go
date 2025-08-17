package order

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUuid string) error {
	status, err := s.orderRepo.GetOrderStatus(ctx, orderUuid)
	if err != nil {
		return err
	}

	if status == model.OrderStatus_PAID {
		return model.ErrOrderAlreadyPaid
	}

	err = s.orderRepo.CancelOrder(ctx, orderUuid)
	if err != nil {
		return err
	}

	return nil
}
