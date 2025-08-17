package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func (s *service) GetPartByID(
	ctx context.Context,
	id string,
) (*model.Part, error) {
	part, err := s.repo.GetByID(ctx, id)
	return part, err
}
