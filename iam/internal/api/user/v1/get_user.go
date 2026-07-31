package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/iam/internal/api/converter"
	"github.com/Reensef/go-microservices-course/iam/internal/model"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func (a *api) GetUser(
	ctx context.Context,
	req *iamV1.GetUserRequest,
) (*iamV1.GetUserResponse, error) {
	user, err := a.service.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		log.Printf("api: error getting user: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrUserUuidInvalidFormat):
			return nil, status.Errorf(codes.InvalidArgument, "user UUID must be UUID format")
		case errors.Is(err, model.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "user with uuid %s not found", req.GetUserUuid())
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	return &iamV1.GetUserResponse{
		User: converter.ToProtoUser(*user),
	}, nil
}
