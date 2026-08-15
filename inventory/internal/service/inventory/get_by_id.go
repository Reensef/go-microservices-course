package service

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
)

func (s *service) GetPartByID(
	ctx context.Context,
	id string,
) (*model.Part, error) {
	ctx, span := tracer.StartSpan(ctx, "inventory.get_part_by_id",
		trace.WithAttributes(attribute.String("part.id", id)),
	)
	defer span.End()

	if len(id) != 24 {
		err := model.ErrPartIdInvalidFormat
		tracer.RecordError(span, err)
		return nil, err
	}

	part, err := s.repo.GetByID(ctx, id)
	if err != nil {
		tracer.RecordError(span, err)
	}

	return part, err
}
