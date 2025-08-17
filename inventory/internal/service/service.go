package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

type InventoryService interface {
	GetPartByID(ctx context.Context, id string) (*model.Part, error)
	GetPartsByFilter(
		ctx context.Context,
		filter *model.PartsFilter,
	) ([]*model.Part, error)
}
