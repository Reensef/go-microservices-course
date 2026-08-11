package service

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/tracing"
)

func (s *service) GetPartByID(
	ctx context.Context,
	id string,
) (*model.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "inventory.get_part_by_id",
		trace.WithAttributes(attribute.String("part.id", id)),
	)
	defer span.End()

	if len(id) != 24 {
		err := model.ErrPartIdInvalidFormat
		tracing.RecordError(span, err)
		return nil, err
	}

	part, err := s.repo.GetByID(ctx, id)
	if err != nil {
		tracing.RecordError(span, err)
	}

	return part, err
}
