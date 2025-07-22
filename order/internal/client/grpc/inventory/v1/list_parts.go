package v1

import (
	"context"

	converter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1/converter"
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (c *inventoryV1Client) ListParts(
	ctx context.Context,
	filter *model.PartsFilter,
) ([]*model.Part, error) {
	response, err := c.service.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: converter.ToProtoFilter(filter),
	})
	if err != nil {
		return nil, err
	}

	protoParts := response.GetParts()

	parts := make([]*model.Part, 0, len(protoParts))
	for _, protoPart := range protoParts {
		part, err := converter.ToModelPart(protoPart)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}
