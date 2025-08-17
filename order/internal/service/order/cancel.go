package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUuid string) error {
	if uuid.Validate(orderUuid) != nil {
		return model.ErrOrderUuidInvalidFormat
	}

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
