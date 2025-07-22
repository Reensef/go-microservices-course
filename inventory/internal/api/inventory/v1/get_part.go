package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
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
	uuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format")
	}

	modelPart, err := a.service.GetPartByUuid(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		log.Printf("api: error getting part by UUID: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.ToProtoPart(modelPart),
	}, nil
}
