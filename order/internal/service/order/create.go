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

	if uuid.Validate(info.UserUuid) != nil {
		return nil, model.ErrUserUuidInvalidFormat
	}

	for _, partId := range info.PartIds {
		if uuid.Validate(partId) != nil {
			return nil, model.ErrPartIdInvalidFormat
		}
	}

	parts, err := s.inventoryService.ListParts(ctx, &model.PartsFilter{
		Ids: info.PartIds,
	})
	if err != nil {
		return nil, err
	}

	partIdsExists := map[string]bool{}
	for _, part := range parts {
		if part == nil {
			continue
		}
		partIdsExists[part.Id] = true
	}

	for _, uuid := range info.PartIds {
		if !partIdsExists[uuid] {
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
