package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) GetOrderByUUID(
	ctx context.Context,
	orderUuid uuid.UUID,
) (*model.Order, error) {
	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUuid)
	return order, err
}
