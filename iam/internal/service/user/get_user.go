package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

func (s *service) GetUser(
	ctx context.Context,
	userUuid string,
) (model.User, error) {
	if uuid.Validate(userUuid) != nil {
		return model.User{}, model.ErrUserUuidInvalidFormat
	}

	return s.userRepo.GetByUUID(ctx, userUuid)
}
