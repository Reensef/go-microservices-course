package v1

import (
	"sync"

	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

var _ grpcClients.IntentoryServiceClient = (*inventoryClient)(nil)

type inventoryClient struct {
	mu      sync.RWMutex
	service inventoryV1.InventoryServiceClient
}

func NewClient(service inventoryV1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{
		service: service,
	}
}
