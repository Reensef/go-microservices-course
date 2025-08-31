package converter

import (
	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToModelPartsFilter(filter *inventoryV1.PartsFilter) *model.PartsFilter {
	if filter == nil {
		return nil
	}

	return &model.PartsFilter{
		IDs:   filter.Ids,
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
	}
}
