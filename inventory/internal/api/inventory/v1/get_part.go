package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1/converter"
	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(
	ctx context.Context,
	req *inventoryV1.GetPartRequest,
) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	modelPart, err := a.service.GetPartByUuid(ctx, &uuid)
	if errors.Is(err, model.ErrPartNotFound) {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.ModelPartToProto(modelPart),
	}, nil
}
