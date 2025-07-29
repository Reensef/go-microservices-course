package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/order/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) CreateOrder(
	ctx context.Context,
	info *model.OrderInfo,
) (*uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[newUUID] = &repoModel.Order{
		Uuid:      newUUID,
		Info:      *repoConverter.ModelOrderInfoToRepoModel(info),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &newUUID, nil
}
