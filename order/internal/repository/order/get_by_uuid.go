package order

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) GetOrderByUUID(
	ctx context.Context,
	orderUuid uuid.UUID,
) (repoModel.Order, error) {
	return repoModel.Order{}, nil
}
