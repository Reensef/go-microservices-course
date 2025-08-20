package order

import (
	"context"
	"fmt"

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
		Ids: info.PartIds,
	})
	if err != nil {
		return nil, err
	}

	partUuidsExists := map[string]bool{}
	for _, part := range parts {
		if part == nil {
			continue
		}
		partUuidsExists[part.Id] = true
	}

	for _, uuid := range info.PartIds {
		if !partUuidsExists[uuid] {
			return nil, model.ErrPartNotFound
		}
	}

	info.TotalPrice = 0
	for _, part := range parts {
		info.TotalPrice += part.Info.Price
	}

	order, err := s.orderRepo.CreateOrder(ctx, info)
	return order, err
}
