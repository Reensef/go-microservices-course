package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func (s *service) GetPartsByFilter(
	ctx context.Context,
	filter *model.PartsFilter,
) ([]*model.Part, error) {
	parts, err := s.repo.GetByFilter(ctx, filter)
	return parts, err
}
