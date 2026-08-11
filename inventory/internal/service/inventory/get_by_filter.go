package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/tracing"
)

func (s *service) GetPartsByFilter(
	ctx context.Context,
	filter *model.PartsFilter,
) ([]*model.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "inventory.get_parts_by_filter")
	defer span.End()

	if filter != nil {
		for _, id := range filter.IDs {
			if len(id) != 24 {
				err := model.ErrPartIdInvalidFormat
				tracing.RecordError(span, err)
				return nil, err
			}
		}
	}

	parts, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		tracing.RecordError(span, err)
	}

	return parts, err
}
