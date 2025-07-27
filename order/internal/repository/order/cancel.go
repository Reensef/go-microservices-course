package order

import (
	"context"

	"github.com/google/uuid"
)

func (r *repository) CancelOrder(
	ctx context.Context,
	orderUuid uuid.UUID,
) error {
	return nil
}
