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
) (model.Session, error) {
	user, passwordHash, err := s.userRepo.GetCredentialsByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return model.Session{}, model.ErrInvalidCredentials
		}

		return model.Session{}, err
	}

	ok, err := hash.Verify(password, passwordHash)
	if err != nil {
		return model.Session{}, err
	}
	if !ok {
		return model.Session{}, model.ErrInvalidCredentials
	}

	return s.sessionRepo.Create(ctx, user.Uuid, s.sessionTTL)
}
