package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) GetOrderByUUID(
	ctx context.Context,
	orderUuid string,
) (*model.Order, error) {
	if uuid.Validate(orderUuid) != nil {
		return nil, model.ErrOrderUuidInvalidFormat
	}

	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUuid)
	return order, err
}
