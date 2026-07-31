package auth

import (
	"time"

	repo "github.com/Reensef/go-microservices-course/iam/internal/repository"
	def "github.com/Reensef/go-microservices-course/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userRepo    repo.UserRepository
	sessionRepo repo.SessionRepository
	sessionTTL  time.Duration
}

func New(userRepo repo.UserRepository, sessionRepo repo.SessionRepository, sessionTTL time.Duration) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}
