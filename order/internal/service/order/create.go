package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func (s *service) CreateOrder(
	ctx context.Context,
	info *model.OrderInfo,
) (*model.Order, error) {
	if info == nil {
		return nil, fmt.Errorf("order info is nil")
	}

	parts, err := s.inventoryService.ListParts(ctx, &model.PartsFilter{
		Uuids: info.PartUuids,
	})
	if err != nil {
		return nil, err
	}

	partUuidsExists := map[uuid.UUID]bool{}
	for _, part := range parts {
		if part == nil {
			continue
		}
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
