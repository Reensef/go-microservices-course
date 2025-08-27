package v1

import (
	"context"

	converter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) GetPart(
	ctx context.Context,
	partId string,
) (*model.Part, error) {
	response, err := c.service.GetPart(
		ctx,
		&inventoryGrpc.GetPartRequest{
			Id: partId,
		},
	)
	if err != nil {
		return nil, err
	}

	return converter.ToModelPart(response.GetPart())
}
