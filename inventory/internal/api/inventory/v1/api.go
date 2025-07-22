package v1

import (
	service "github.com/Reensef/go-microservices-course/inventory/internal/service"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	service service.InventoryService
}

func NewAPI(service service.InventoryService) *api {
	return &api{
		service: service,
	}
}
