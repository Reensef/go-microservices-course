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

func (a *api) Whoami(
	ctx context.Context,
	req *iamV1.WhoamiRequest,
) (*iamV1.WhoamiResponse, error) {
	session, user, err := a.service.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		log.Printf("api: error getting whoami: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrSessionUuidInvalidFormat):
			return nil, status.Errorf(codes.InvalidArgument, "session UUID must be UUID format")
		case errors.Is(err, model.ErrSessionNotFound):
			return nil, status.Errorf(codes.NotFound, "session not found")
		case errors.Is(err, model.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "user not found")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	return &iamV1.WhoamiResponse{
		Session: converter.ToProtoSession(*session),
		User:    converter.ToProtoUser(*user),
	}, nil
}
