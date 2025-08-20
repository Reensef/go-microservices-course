package repository

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

type PartRepository interface {
	Create(
		ctx context.Context,
		part *model.PartInfo,
	) (*model.Part, error)

	GetByID(
		ctx context.Context,
		id string,
	) (*model.Part, error)

	GetAll(
		ctx context.Context,
	) ([]*model.Part, error)

	GetByFilter(
		ctx context.Context,
		filter *model.PartsFilter,
	) ([]*model.Part, error)
}
