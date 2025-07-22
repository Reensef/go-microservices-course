package order

import (
	"context"

	"github.com/google/uuid"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/order/internal/repository/converter"
)

func (r *repository) GetOrderByUUID(
	ctx context.Context,
	orderUuid uuid.UUID,
) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[orderUuid]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return repoConverter.ToModelOrder(order), nil
}
