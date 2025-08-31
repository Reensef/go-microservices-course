package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1/converter"
	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryV1.ListPartsRequest,
) (*inventoryV1.ListPartsResponse, error) {
	modelFilter := converter.ToModelPartsFilter(req.GetFilter())

	modelParts, err := a.service.GetPartsByFilter(ctx, modelFilter)
	if err != nil {
		log.Printf("api: error getting parts by filter: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrPartIdInvalidFormat):
			return nil, status.Errorf(codes.InvalidArgument, "part id must be ObjectID format")
		case errors.Is(err, model.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "part not found")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	protoParts := converter.ToProtoParts(modelParts)

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
