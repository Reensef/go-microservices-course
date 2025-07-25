package order

import (
	def "github.com/Reensef/go-microservices-course/order/internal/service"
	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	ufoRepository repo.OrderRepository
}

func NewService(ufoRepository repo.OrderRepository) *service {
	return &service{
		ufoRepository: ufoRepository,
	}
}
