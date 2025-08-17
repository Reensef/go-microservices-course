package service

import (
	repo "github.com/Reensef/go-microservices-course/inventory/internal/repository"
)

type service struct {
	repo repo.PartRepository
}

func NewService(repo repo.PartRepository) *service {
	service := &service{
		repo: repo,
	}

	return service
}
