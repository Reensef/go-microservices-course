package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

type AuthService interface {
	Login(
		ctx context.Context,
		login string,
		password string,
	) (model.Session, error)

	Whoami(
		ctx context.Context,
		sessionUuid string,
	) (model.Session, model.User, error)
}

type UserService interface {
	Register(
		ctx context.Context,
		info *model.UserRegistrationInfo,
	) (model.User, error)

	GetUser(
		ctx context.Context,
		userUuid string,
	) (model.User, error)
}
