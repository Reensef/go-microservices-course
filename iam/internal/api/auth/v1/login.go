package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func (a *api) Login(
	ctx context.Context,
	req *iamV1.LoginRequest,
) (*iamV1.LoginResponse, error) {
	session, err := a.service.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		logger.Error("api: error logging in", zap.Error(err))

		switch {
		case errors.Is(err, model.ErrInvalidCredentials):
			return nil, status.Errorf(codes.Unauthenticated, "invalid login or password")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	logger.Info("user logged in", zap.String("user_uuid", session.UserUuid))

	return &iamV1.LoginResponse{
		SessionUuid: session.Uuid,
	}, nil
}
