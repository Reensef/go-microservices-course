package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

type InventoryService interface {
	GetPartByUuid(ctx context.Context, uuid uuid.UUID) (*model.Part, error)
	GetPartsByFilter(
		ctx context.Context,
		filter *model.PartsFilter,
	) []*model.Part
}
