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

func (a *api) GetPart(
	ctx context.Context,
	req *inventoryV1.GetPartRequest,
) (*inventoryV1.GetPartResponse, error) {
	modelPart, err := a.service.GetPartByID(ctx, req.GetId())
	if err != nil {
		log.Printf("api: error getting part by UUID: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "part with id %s not found", req.GetId())
		case errors.Is(err, model.ErrPartIdInvalidFormat):
			return nil, status.Errorf(codes.InvalidArgument, "part id must be ObjectID format")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.ToProtoPart(modelPart),
	}, nil
}
