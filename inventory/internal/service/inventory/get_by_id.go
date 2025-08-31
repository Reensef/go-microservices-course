package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func (s *service) GetPartByID(
	ctx context.Context,
	id string,
) (*model.Part, error) {
	if len(id) != 24 {
		return nil, model.ErrPartIdInvalidFormat
	}

	part, err := s.repo.GetByID(ctx, id)
	return part, err
}
