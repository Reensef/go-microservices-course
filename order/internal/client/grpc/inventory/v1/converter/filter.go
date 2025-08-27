package converter

import (
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToProtoFilter(filter *model.PartsFilter) *inventoryGrpc.PartsFilter {
	categories := make([]inventoryGrpc.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, ToProtoPartCategory(category))
	}

	return &inventoryGrpc.PartsFilter{
		Ids:                   filter.Ids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
