package v1

import (
	"context"
	"errors"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/iam/internal/api/converter"
	"github.com/Reensef/go-microservices-course/iam/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func (a *api) Register(
	ctx context.Context,
	req *iamV1.RegisterRequest,
) (*iamV1.RegisterResponse, error) {
	registrationInfo := &model.UserRegistrationInfo{
		Info:     converter.ToModelUserInfo(req.GetInfo().GetInfo()),
		Password: req.GetInfo().GetPassword(),
	}

	user, err := a.service.Register(ctx, registrationInfo)
	if err != nil {
		log.Printf("api: error registering user: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrUserAlreadyExists):
			return nil, status.Errorf(codes.AlreadyExists, "user with this login already exists")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	logger.Info("user registered", zap.String("user_uuid", user.Uuid))

	return &iamV1.RegisterResponse{
		UserUuid: user.Uuid,
	}, nil
}
