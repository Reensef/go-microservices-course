package order

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
	def "github.com/Reensef/go-microservices-course/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo        repo.OrderRepository
	inventoryService grpcClients.IntentoryServiceClient
	paymentService   grpcClients.PaymentServiceClient
}

func NewService(
	orderRepo repo.OrderRepository,
	inventoryService grpcClients.IntentoryServiceClient,
	paymentService grpcClients.PaymentServiceClient,
) *service {
	return &service{
		orderRepo:        orderRepo,
		inventoryService: inventoryService,
		paymentService:   paymentService,
	}
}
