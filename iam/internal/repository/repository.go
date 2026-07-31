package repository

import (
	"context"
	"time"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		info *model.UserRegistrationInfo,
		passwordHash string,
	) (*model.User, error)

	GetByUUID(
		ctx context.Context,
		userUuid string,
	) (*model.User, error)

	// GetCredentialsByLogin возвращает пользователя вместе с хешем пароля для проверки при входе
	GetCredentialsByLogin(
		ctx context.Context,
		login string,
	) (*model.User, string, error)
}

type SessionRepository interface {
	Create(
		ctx context.Context,
		userUuid string,
		ttl time.Duration,
	) (*model.Session, error)

	GetByUUID(
		ctx context.Context,
		sessionUuid string,
	) (*model.Session, error)
}
