package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1/converter"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryV1.ListPartsRequest,
) (*inventoryV1.ListPartsResponse, error) {
	modelFilter, err := converter.ToModelPartsFilter(req.GetFilter())
	if err != nil {
		if errors.Is(err, converter.ErrParseUuid) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format")
		}
		log.Printf("api: error parsing filter: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	modelParts := a.service.GetPartsByFilter(ctx, modelFilter)
	protoParts := converter.ToProtoParts(modelParts)

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
