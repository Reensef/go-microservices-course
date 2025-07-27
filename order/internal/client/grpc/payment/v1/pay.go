package v1

import (
	"context"

	"github.com/google/uuid"
)

func (c *paymentClient) PayOrder(
	ctx context.Context,
	orderUuid, userUuid *uuid.UUID,
) (*uuid.UUID, error) {
	return nil, nil
}
