package order

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	kafka "github.com/Reensef/go-microservices-course/order/internal/events"
	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
	def "github.com/Reensef/go-microservices-course/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo        repo.OrderRepository
	inventoryService grpcClients.IntentoryClient
	paymentService   grpcClients.PaymentClient
	orderProducer    kafka.OrderProducer
}

func New(
	orderRepo repo.OrderRepository,
	inventoryService grpcClients.IntentoryClient,
	paymentService grpcClients.PaymentClient,
	orderProducer kafka.OrderProducer,
) *service {
	return &service{
		orderRepo:        orderRepo,
		inventoryService: inventoryService,
		paymentService:   paymentService,
	}
}
