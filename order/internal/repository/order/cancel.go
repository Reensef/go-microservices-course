package order

import (
	"context"

	"github.com/google/uuid"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) CancelOrder(
	ctx context.Context,
	orderUuid uuid.UUID,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[orderUuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	if order.Info.Status == repoModel.OrderStatus_PAID {
		return model.ErrOrderAlreadyPaid
	}

	order.Info.Status = repoModel.OrderStatus_CANCELED

	return nil
}
