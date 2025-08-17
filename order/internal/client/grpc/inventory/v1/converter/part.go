package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

func ToModelPart(part *inventoryV1.Part) (*model.Part, error) {
	modelPart := &model.Part{
		Id: part.GetId(),
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

	return modelPart, nil
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

func ToModelPartDimensions(dim *inventoryV1.Dimensions) *model.PartDimensions {
	return &model.PartDimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ToModelPartManufacturer(manufacturer *inventoryV1.Manufacturer) *model.PartManufacturer {
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToProtoPart(part *model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Id:            part.Id,
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

func ToProtoPartDimensions(dim *model.PartDimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ToProtoPartManufacturer(manufacturer *model.PartManufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
