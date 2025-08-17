package converter

import (
	"maps"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func ToRepoModelPart(part *model.Part) *repoModel.Part {
	if part == nil {
		return nil
	}
	return &repoModel.Part{
		ID:        part.ID,
		Info:      *ToRepoModelPartInfo(&part.Info),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func ToRepoModelPartInfo(info *model.PartInfo) *repoModel.PartInfo {
	if info == nil {
		return nil
	}
	return &repoModel.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      ToRepoModelPartCategory(info.Category),
		Dimensions:    ToRepoModelPartDimensions(info.Dimensions),
		Manufacturer:  ToRepoModelPartManufacturer(info.Manufacturer),
		Tags:          info.Tags,
		Metadata:      maps.Clone(info.Metadata),
	}
}

func ToRepoModelPartCategory(category model.PartCategory) repoModel.PartCategory {
	switch category {
	case model.PartCategory_ENGINE:
		return repoModel.PartCategory_ENGINE
	case model.PartCategory_FUEL:
		return repoModel.PartCategory_FUEL
	case model.PartCategory_PORTHOLE:
		return repoModel.PartCategory_PORTHOLE
	case model.PartCategory_WING:
		return repoModel.PartCategory_WING
	default:
		return repoModel.PartCategory_UNSPECIFIED
	}
}

func ToRepoModelPartManufacturer(manufacturer *model.PartManufacturer) *repoModel.PartManufacturer {
	if manufacturer == nil {
		return nil
	}
	return &repoModel.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToRepoModelPartDimensions(dim *model.PartDimensions) *repoModel.PartDimensions {
	if dim == nil {
		return nil
	}
	return &repoModel.PartDimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ToModelPart(part *repoModel.Part) *model.Part {
	if part == nil {
		return nil
	}
	return &model.Part{
		ID:        part.ID,
		Info:      *ToModelPartInfo(&part.Info),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func ToModelParts(repoParts []*repoModel.Part) []*model.Part {
	modelParts := make([]*model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		modelParts = append(modelParts, ToModelPart(repoPart))
	}
	return modelParts
}

func ToModelPartInfo(info *repoModel.PartInfo) *model.PartInfo {
	if info == nil {
		return nil
	}
	return &model.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      ToModelPartCategory(info.Category),
		Dimensions:    ToModelPartDimensions(info.Dimensions),
		Manufacturer:  ToModelPartManufacturer(info.Manufacturer),
		Tags:          info.Tags,
		Metadata:      maps.Clone(info.Metadata),
	}
}

func ToModelPartCategory(category repoModel.PartCategory) model.PartCategory {
	switch category {
	case repoModel.PartCategory_ENGINE:
		return model.PartCategory_ENGINE
	case repoModel.PartCategory_FUEL:
		return model.PartCategory_FUEL
	case repoModel.PartCategory_PORTHOLE:
		return model.PartCategory_PORTHOLE
	case repoModel.PartCategory_WING:
		return model.PartCategory_WING
	default:
		return model.PartCategory_UNSPECIFIED
	}
}

func ToModelPartDimensions(dim *repoModel.PartDimensions) *model.PartDimensions {
	if dim == nil {
		return nil
	}
	return &model.PartDimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ToModelPartManufacturer(manufacturer *repoModel.PartManufacturer) *model.PartManufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
