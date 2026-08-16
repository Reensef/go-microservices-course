package user

import (
	"context"

	"github.com/Reensef/go-microservices-course/iam/internal/hash"
	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

func (s *service) Register(
	ctx context.Context,
	info *model.UserRegistrationInfo,
) (model.User, error) {
	passwordHash, err := hash.Hash(info.Password)
	if err != nil {
		return model.User{}, err
	}

	return s.userRepo.Create(ctx, info, passwordHash)
}
