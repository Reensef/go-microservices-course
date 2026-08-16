package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

func (s *service) Whoami(
	ctx context.Context,
	sessionUuid string,
) (model.Session, model.User, error) {
	if uuid.Validate(sessionUuid) != nil {
		return model.Session{}, model.User{}, model.ErrSessionUuidInvalidFormat
	}

	session, err := s.sessionRepo.GetByUUID(ctx, sessionUuid)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UserUuid)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	return session, user, nil
}
