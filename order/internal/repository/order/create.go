package order

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) CreateOrder(
	ctx context.Context,
	userUuid uuid.UUID,
	partUuids []uuid.UUID,
) (repoModel.Order, error) {
	return repoModel.Order{}, nil
}
