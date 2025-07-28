package order

import (
	"context"

	"github.com/google/uuid"

	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) CreateOrder(
	ctx context.Context,
	info *model.OrderInfo,
) (*uuid.UUID, error) {
	parts, err := s.inventoryService.ListParts(ctx, &grpcClients.PartsFilter{
		Uuids: info.PartUuids,
	})
	if err != nil {
		return nil, err
	}

	partUuidsExists := map[uuid.UUID]bool{}
	for _, part := range parts {
		partUuidsExists[part.Uuid] = true
	}

	for _, uuid := range info.PartUuids {
		if !partUuidsExists[uuid] {
			return nil, model.ErrPartNotFound
		}
	}

	info.TotalPrice = 0
	for _, part := range parts {
		info.TotalPrice += part.Info.Price
	}

	return s.orderRepo.CreateOrder(ctx, info)
}
