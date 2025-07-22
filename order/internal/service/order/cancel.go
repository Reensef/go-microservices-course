package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) CancelOrder(ctx context.Context, orderUuid uuid.UUID) error {
	err := s.orderRepo.CancelOrder(ctx, orderUuid)
	return err
}
