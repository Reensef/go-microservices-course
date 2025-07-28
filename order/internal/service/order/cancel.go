package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) CancelOrder(ctx context.Context, orderUuid *uuid.UUID) error {
	return s.orderRepo.CancelOrder(ctx, orderUuid)
}
