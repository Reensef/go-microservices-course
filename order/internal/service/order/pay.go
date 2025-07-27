package order

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) PayOrder(
	ctx context.Context,
	id uuid.UUID,
	paymentMethod model.OrderPaymentMethod,
) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
