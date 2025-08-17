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

	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUuid)
	if err != nil {
		
	}

	if order.Info.Status == model.OrderStatus_PAID {
		return model.ErrOrderAlreadyPaid
	}

	err = s.orderRepo.CancelOrder(ctx, orderUuid)
	if err != nil {
		return err
	}

	return nil
}
