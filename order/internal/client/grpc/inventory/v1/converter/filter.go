package converter

import (
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToProtoFilter(filter *model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, ToProtoPartCategory(category))
	}

	uuidStrings := make([]string, 0, len(filter.Uuids))
	for _, uuid := range filter.Uuids {
		uuidStrings = append(uuidStrings, uuid.String())
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuidStrings,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
