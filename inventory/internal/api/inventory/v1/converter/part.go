package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

func ModelPartToProto(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return nil
	}
	return &inventoryV1.Part{
		Uuid:          part.Uuid.String(),
		Name:          part.Info.Name,
		Description:   part.Info.Description,
		Price:         part.Info.Price,
		StockQuantity: part.Info.StockQuantity,
		Category:      ModelPartCategoryToProto(part.Info.Category),
		Dimensions:    ModelPartDimensionsToProto(part.Info.Dimensions),
		Manufacturer:  ModelPartManufacturerToProto(part.Info.Manufacturer),
		Tags:          part.Info.Tags,
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
		Metadata: func() map[string]*inventoryV1.Value {
			metadata := make(map[string]*inventoryV1.Value)
			for key, value := range part.Info.Metadata {
				metadata[key] = MultiValueToProto(&value)
			}
			return metadata
		}(),
	}
}

func ModelPartsToProto(parts []*model.Part) []*inventoryV1.Part {
	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, ModelPartToProto(part))
	}
	return result
}

func ProtoPartToModel(part *inventoryV1.Part) (*model.Part, error) {
	if part == nil {
		return nil, nil
	}
	uuid, err := uuid.Parse(part.Uuid)
	if err != nil {
		return nil, err
	}

	return &model.Part{
		Uuid: uuid,
		Info: model.PartInfo{
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      ProtoPartCategoryToModel(part.Category),
			Dimensions:    ProtoPartDimensionsToModel(part.Dimensions),
			Manufacturer:  ProtoPartManufacturerToModel(part.Manufacturer),
			Tags:          part.Tags,
			Metadata: func() map[string]utils.MultiValue {
				metadata := make(map[string]utils.MultiValue)
				for key, value := range part.Metadata {
					metadata[key] = *ProtoValueToMultiValue(value)
				}
				return metadata
			}(),
		},
		CreatedAt: part.CreatedAt.AsTime(),
		UpdatedAt: part.UpdatedAt.AsTime(),
	}, nil
}

func ModelPartCategoryToProto(category model.PartCategory) inventoryV1.Category {
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

func ProtoPartCategoryToModel(category inventoryV1.Category) model.PartCategory {
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

func ModelPartDimensionsToProto(dimensions *model.PartDimensions) *inventoryV1.Dimensions {
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

func ProtoPartDimensionsToModel(dimensions *inventoryV1.Dimensions) *model.PartDimensions {
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

func ModelPartManufacturerToProto(manufacturer *model.PartManufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ProtoPartManufacturerToModel(manufacturer *inventoryV1.Manufacturer) *model.PartManufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
