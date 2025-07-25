package order

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) CreateOrder(
	ctx context.Context,
	userUuid uuid.UUID,
	partUuids []uuid.UUID,
) (model.Order, error) {
	return model.Order{}, nil
}

func (s *service) GetOrderByUUID(
	ctx context.Context,
	orderUuid uuid.UUID,
) (model.Order, error) {
	return model.Order{}, nil
}

func (s *service) CancelOrder(ctx context.Context, orderUuid uuid.UUID) error {
	return nil
}

func (s *service) PayOrder(
	ctx context.Context,
	id uuid.UUID,
	paymentMethod model.OrderPaymentMethod,
) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
