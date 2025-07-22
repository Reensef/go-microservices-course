package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

type PartRepository interface {
	GetByUuid(
		ctx context.Context,
		uuid uuid.UUID,
	) (*model.Part, error)

	GetByFilter(
		ctx context.Context,
		filter *model.PartsFilter,
	) []*model.Part
}
