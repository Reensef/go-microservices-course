package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/multivalue"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToModelPart(part *inventoryGrpc.Part) (*model.Part, error) {
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
			Metadata: func() map[string]multivalue.MultiValue {
				metadata := make(map[string]multivalue.MultiValue)
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

func ToModelPartCategory(category inventoryGrpc.Category) model.PartCategory {
	switch category {
	case inventoryGrpc.Category_CATEGORY_ENGINE:
		return model.PartCategory_ENGINE
	case inventoryGrpc.Category_CATEGORY_FUEL:
		return model.PartCategory_FUEL
	case inventoryGrpc.Category_CATEGORY_PORTHOLE:
		return model.PartCategory_PORTHOLE
	case inventoryGrpc.Category_CATEGORY_WING:
		return model.PartCategory_WING
	default:
		return model.PartCategory_UNSPECIFIED
	}
}

func ToModelPartDimensions(dim *inventoryGrpc.Dimensions) *model.PartDimensions {
	return &model.PartDimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ToModelPartManufacturer(manufacturer *inventoryGrpc.Manufacturer) *model.PartManufacturer {
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToProtoPart(part *model.Part) *inventoryGrpc.Part {
	return &inventoryGrpc.Part{
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
		Metadata: func() map[string]*inventoryGrpc.Value {
			metadata := make(map[string]*inventoryGrpc.Value)
			for key, value := range part.Info.Metadata {
				metadata[key] = ToProtoMultiValue(&value)
			}
			return metadata
		}(),
	}
}

func ToProtoPartCategory(category model.PartCategory) inventoryGrpc.Category {
	switch category {
	case model.PartCategory_ENGINE:
		return inventoryGrpc.Category_CATEGORY_ENGINE
	case model.PartCategory_FUEL:
		return inventoryGrpc.Category_CATEGORY_FUEL
	case model.PartCategory_PORTHOLE:
		return inventoryGrpc.Category_CATEGORY_PORTHOLE
	case model.PartCategory_WING:
		return inventoryGrpc.Category_CATEGORY_WING
	default:
		return inventoryGrpc.Category_CATEGORY_UNSPECIFIED
	}
}

func ToProtoPartDimensions(dim *model.PartDimensions) *inventoryGrpc.Dimensions {
	return &inventoryGrpc.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ToProtoPartManufacturer(manufacturer *model.PartManufacturer) *inventoryGrpc.Manufacturer {
	return &inventoryGrpc.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
