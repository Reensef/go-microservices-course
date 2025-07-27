package order

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) PayOrder(
	ctx context.Context,
	orderUuid uuid.UUID,
	paymentMethod repoModel.OrderPaymentMethod,
) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
