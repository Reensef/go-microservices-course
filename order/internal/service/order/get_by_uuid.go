package order

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) GetOrderByUUID(
	ctx context.Context,
	orderUuid uuid.UUID,
) (model.Order, error) {
	return model.Order{}, nil
}
