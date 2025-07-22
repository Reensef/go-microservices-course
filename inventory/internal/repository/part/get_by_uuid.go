package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	converter "github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
)

func (r *repository) GetByUuid(
	ctx context.Context,
	uuid uuid.UUID,
) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return converter.ToModelPart(order), nil
}
