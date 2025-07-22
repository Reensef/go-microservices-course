package v1

import (
	"context"

	"github.com/google/uuid"

	converter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (c *inventoryV1Client) GetPart(
	ctx context.Context,
	partUuid uuid.UUID,
) (*model.Part, error) {
	response, err := c.service.GetPart(
		ctx,
		&inventoryV1.GetPartRequest{
			Uuid: partUuid.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return converter.ToModelPart(response.GetPart())
}
