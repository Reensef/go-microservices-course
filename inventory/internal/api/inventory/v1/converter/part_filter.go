package converter

import (
	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToModelPartsFilter(filter *inventoryV1.PartsFilter) (*model.PartsFilter, error) {
	if filter == nil {
		return nil, nil
	}
	uuids := make([]uuid.UUID, 0, len(filter.GetUuids()))
	for _, uuidString := range filter.GetUuids() {
		uuid, err := uuid.Parse(uuidString)
		if err != nil {
			return nil, ErrParseUuid
		}
		uuids = append(uuids, uuid)
	}

	return &model.PartsFilter{
		Uuids: uuids,
		Names: filter.Names,
		Categories: func() []model.PartCategory {
			categories := make([]model.PartCategory, 0, len(filter.Categories))
			for _, category := range filter.Categories {
				categories = append(categories, ToModelPartCategory(category))
			}
			return categories
		}(),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}, nil
}
