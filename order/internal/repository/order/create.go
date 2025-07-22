package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoConverter "github.com/Reensef/go-microservices-course/order/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func (r *repository) CreateOrder(
	ctx context.Context,
	info *model.OrderInfo,
) (*model.Order, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("order info is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := &repoModel.Order{
		Uuid:      newUUID,
		Info:      *repoConverter.ToRepoModelOrderInfo(info),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.data[newUUID] = repoOrder

	return &model.Order{
		Uuid:      repoOrder.Uuid,
		Info:      *info,
		CreatedAt: repoOrder.CreatedAt,
		UpdatedAt: repoOrder.UpdatedAt,
	}, nil
}
