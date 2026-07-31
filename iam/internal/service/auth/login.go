package auth

import (
	"context"
	"errors"

	"github.com/Reensef/go-microservices-course/iam/internal/hash"
	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

func (s *service) Login(
	ctx context.Context,
	login string,
	password string,
) (*model.Session, error) {
	user, passwordHash, err := s.userRepo.GetCredentialsByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidCredentials
		}

		return nil, err
	}

	ok, err := hash.Verify(password, passwordHash)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	return s.sessionRepo.Create(ctx, user.Uuid, s.sessionTTL)
}
