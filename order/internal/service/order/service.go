package order

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
	def "github.com/Reensef/go-microservices-course/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo        repo.OrderRepository
	inventoryService grpcClients.IntentoryClient
	paymentService   grpcClients.PaymentClient
}

func New(
	orderRepo repo.OrderRepository,
	inventoryService grpcClients.IntentoryClient,
	paymentService grpcClients.PaymentClient,
) *service {
	return &service{
		orderRepo:        orderRepo,
		inventoryService: inventoryService,
		paymentService:   paymentService,
	}
}
