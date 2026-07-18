package v1

import (
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

var _ grpcClients.IntentoryClient = (*inventoryClient)(nil)

type inventoryClient struct {
	service inventoryGrpc.InventoryServiceClient
}

func New(service inventoryGrpc.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{
		service: service,
	}
}
