package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func (s *service) GetPartsByFilter(
	ctx context.Context,
	filter *model.PartsFilter,
) ([]*model.Part, error) {
	if filter != nil {
		for _, id := range filter.IDs {
			if len(id) != 24 {
				return nil, model.ErrPartIdInvalidFormat
			}
		}
	}

	parts, err := s.repo.GetByFilter(ctx, filter)
	return parts, err
}
