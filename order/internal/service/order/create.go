package order

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

// objectIDHexLen — длина hex-строки MongoDB ObjectID (12 байт = 24 hex-символа).
const objectIDHexLen = 24

func isValidObjectIDHex(id string) bool {
	if len(id) != objectIDHexLen {
		return false
	}
	_, err := hex.DecodeString(id)
	return err == nil
}

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
		if !isValidObjectIDHex(partId) {
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
