package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1/converter"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryV1.ListPartsRequest,
) (*inventoryV1.ListPartsResponse, error) {
	modelFilter, err := converter.ProtoFilterToModel(req.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	modelParts := a.service.GetPartsByFilter(ctx, modelFilter)

	protoParts := make([]*inventoryV1.Part, 0, len(modelParts))
	for _, modelPart := range modelParts {
		protoParts = append(protoParts, converter.ModelPartToProto(modelPart))
	}

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
