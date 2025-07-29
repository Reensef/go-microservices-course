package converter

import (
	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func ModelPartToRepoModel(part *model.Part) *repoModel.Part {
	return &repoModel.Part{
		Uuid:      part.Uuid,
		Info:      *ModelPartInfoToRepoModel(&part.Info),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func RepoModelPartToModel(part *repoModel.Part) *model.Part {
	return &model.Part{
		Uuid:      part.Uuid,
		Info:      *RepoModelPartInfoToModel(&part.Info),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func ModelPartInfoToRepoModel(info *model.PartInfo) *repoModel.PartInfo {
	return &repoModel.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      ModelPartCategoryToRepoModel(info.Category),
		Dimensions:    ModelPartDimensionsToRepoModel(info.Dimensions),
		Manufacturer:  ModelPartManufacturerToRepoModel(info.Manufacturer),
		Tags:          info.Tags,
		Metadata:      info.Metadata,
	}
}

func RepoModelPartInfoToModel(info *repoModel.PartInfo) *model.PartInfo {
	return &model.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      RepoModelPartCategoryToModel(info.Category),
		Dimensions:    RepoModelPartDimensionsToModel(info.Dimensions),
		Manufacturer:  RepoModelPartManufacturerToModel(info.Manufacturer),
		Tags:          info.Tags,
		Metadata:      info.Metadata,
	}
}

func ModelPartCategoryToRepoModel(category model.PartCategory) repoModel.PartCategory {
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

func RepoModelPartCategoryToModel(category repoModel.PartCategory) model.PartCategory {
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

func ModelPartDimensionsToRepoModel(dim *model.PartDimensions) *repoModel.PartDimensions {
	return &repoModel.PartDimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func RepoModelPartDimensionsToModel(dim *repoModel.PartDimensions) *model.PartDimensions {
	return &model.PartDimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ModelPartManufacturerToRepoModel(manufacturer *model.PartManufacturer) *repoModel.PartManufacturer {
	return &repoModel.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func RepoModelPartManufacturerToModel(manufacturer *repoModel.PartManufacturer) *model.PartManufacturer {
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
