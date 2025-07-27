package v1

import (
	"context"

	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	converter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1/converter"
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(
	ctx context.Context,
	filter *grpcClients.PartsFilter,
) ([]*model.Part, error) {
	response, err := c.service.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: converter.ModelFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	parts := make([]*model.Part, 0, len(response.Parts))
	for _, protoPart := range response.Parts {
		part, err := converter.ProtoPartToModel(protoPart)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}
