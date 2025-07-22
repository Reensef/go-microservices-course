package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func (s *service) GetPartByUuid(
	ctx context.Context,
	uuid uuid.UUID,
) (*model.Part, error) {
	part, err := s.repo.GetByUuid(ctx, uuid)
	return part, err
}
