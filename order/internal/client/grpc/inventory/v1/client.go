package v1

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

var _ grpcClients.IntentoryServiceClient = (*inventoryV1Client)(nil)

type inventoryV1Client struct {
	service inventoryV1.InventoryServiceClient
}

func NewClient(service inventoryV1.InventoryServiceClient) *inventoryV1Client {
	return &inventoryV1Client{
		service: service,
	}
}
