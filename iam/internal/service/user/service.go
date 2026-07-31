package user

import (
	repo "github.com/Reensef/go-microservices-course/iam/internal/repository"
	def "github.com/Reensef/go-microservices-course/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepo repo.UserRepository
}

func New(userRepo repo.UserRepository) *service {
	return &service{
		userRepo: userRepo,
	}
}
