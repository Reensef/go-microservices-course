package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func (s *service) GetPartsByFilter(
	ctx context.Context,
	filter *model.PartsFilter,
) []*model.Part {
	parts := s.repo.GetByFilter(ctx, filter)
	return parts
}
