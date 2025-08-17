package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

func ToProtoPart(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return nil
	}
	return &inventoryV1.Part{
		Id:            part.ID,
		Name:          part.Info.Name,
		Description:   part.Info.Description,
		Price:         part.Info.Price,
		StockQuantity: part.Info.StockQuantity,
		Category:      ToProtoPartCategory(part.Info.Category),
		Dimensions:    ToProtoPartDimensions(part.Info.Dimensions),
		Manufacturer:  ToProtoPartManufacturer(part.Info.Manufacturer),
		Tags:          part.Info.Tags,
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
		Metadata: func() map[string]*inventoryV1.Value {
			metadata := make(map[string]*inventoryV1.Value)
			for key, value := range part.Info.Metadata {
				metadata[key] = ToProtoMultiValue(&value)
			}
			return metadata
		}(),
	}
}

func ToProtoParts(parts []*model.Part) []*inventoryV1.Part {
	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, ToProtoPart(part))
	}
	return result
}

func ToModelPart(part *inventoryV1.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		ID: part.Id,
		Info: model.PartInfo{
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      ToModelPartCategory(part.Category),
			Dimensions:    ToModelPartDimensions(part.Dimensions),
			Manufacturer:  ToModelPartManufacturer(part.Manufacturer),
			Tags:          part.Tags,
			Metadata: func() map[string]utils.MultiValue {
				metadata := make(map[string]utils.MultiValue)
				for key, value := range part.Metadata {
					metadata[key] = *ToMultiValue(value)
				}
				return metadata
			}(),
		},
		CreatedAt: part.CreatedAt.AsTime(),
		UpdatedAt: part.UpdatedAt.AsTime(),
	}
}

func ToProtoPartCategory(category model.PartCategory) inventoryV1.Category {
	switch category {
	case model.PartCategory_ENGINE:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.PartCategory_FUEL:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.PartCategory_PORTHOLE:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.PartCategory_WING:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func ToModelPartCategory(category inventoryV1.Category) model.PartCategory {
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.PartCategory_ENGINE
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.PartCategory_FUEL
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.PartCategory_PORTHOLE
	case inventoryV1.Category_CATEGORY_WING:
		return model.PartCategory_WING
	default:
		return model.PartCategory_UNSPECIFIED
	}
}

func ToProtoPartDimensions(dimensions *model.PartDimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToModelPartDimensions(dimensions *inventoryV1.Dimensions) *model.PartDimensions {
	if dimensions == nil {
		return nil
	}
	return &model.PartDimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToProtoPartManufacturer(manufacturer *model.PartManufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToModelPartManufacturer(manufacturer *inventoryV1.Manufacturer) *model.PartManufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
