package converter

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func TestModelPartToRepoModel(t *testing.T) {
	part := &model.Part{
		Uuid:      uuid.MustParse(gofakeit.UUID()),
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
		Info: model.PartInfo{
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(10, 1000),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      model.PartCategory(gofakeit.Number(0, 4)),
			Dimensions: &model.PartDimensions{
				Length: gofakeit.Float64Range(1, 100),
				Width:  gofakeit.Float64Range(1, 100),
				Height: gofakeit.Float64Range(1, 100),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: &model.PartManufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: func() []string {
				tags := []string{}
				for i := 0; i < gofakeit.Number(1, 5); i++ {
					tags = append(tags, gofakeit.Word())
				}
				return tags
			}(),
			Metadata: func() map[string]*any {
				meta := map[string]*any{}
				for i := 0; i < gofakeit.Number(1, 5); i++ {
					value := any(gofakeit.Word())
					meta[gofakeit.Word()] = &value
				}
				return meta
			}(),
		},
	}

	repoPart := ModelPartToRepoModel(part)

	assert.Equal(t, part.Uuid, repoPart.Uuid)
	assert.Equal(t, part.Info.Name, repoPart.Info.Name)
	assert.Equal(t, part.Info.Description, repoPart.Info.Description)
	assert.Equal(t, part.Info.Price, repoPart.Info.Price)
	assert.Equal(t, part.Info.StockQuantity, repoPart.Info.StockQuantity)
	assert.Equal(t, ModelPartCategoryToRepoModel(part.Info.Category), repoPart.Info.Category)
	assert.Equal(t, part.Info.Tags, repoPart.Info.Tags)
	assert.Equal(t, part.Info.Metadata, repoPart.Info.Metadata)
	assert.Equal(t, part.CreatedAt, repoPart.CreatedAt)
	assert.Equal(t, part.UpdatedAt, repoPart.UpdatedAt)
}

func TestRepoModelPartToModel(t *testing.T) {
	repoPart := &repoModel.Part{
		Uuid:      uuid.MustParse(gofakeit.UUID()),
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
		Info: repoModel.PartInfo{
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(10, 1000),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      repoModel.PartCategory(gofakeit.Number(0, 4)),
			Dimensions: &repoModel.PartDimensions{
				Length: gofakeit.Float64Range(1, 100),
				Width:  gofakeit.Float64Range(1, 100),
				Height: gofakeit.Float64Range(1, 100),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: &repoModel.PartManufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: func() []string {
				tags := []string{}
				for i := 0; i < gofakeit.Number(1, 5); i++ {
					tags = append(tags, gofakeit.Word())
				}
				return tags
			}(),
			Metadata: func() map[string]*any {
				meta := map[string]*any{}
				for i := 0; i < gofakeit.Number(1, 5); i++ {
					value := any(gofakeit.Word())
					meta[gofakeit.Word()] = &value
				}
				return meta
			}(),
		},
	}

	part := RepoModelPartToModel(repoPart)

	assert.Equal(t, repoPart.Uuid, part.Uuid)
	assert.Equal(t, repoPart.Info.Name, part.Info.Name)
	assert.Equal(t, repoPart.Info.Description, part.Info.Description)
	assert.Equal(t, repoPart.Info.Price, part.Info.Price)
	assert.Equal(t, repoPart.Info.StockQuantity, part.Info.StockQuantity)
	assert.Equal(t, RepoModelPartCategoryToModel(repoPart.Info.Category), part.Info.Category)
	assert.Equal(t, repoPart.Info.Tags, part.Info.Tags)
	assert.Equal(t, repoPart.Info.Metadata, part.Info.Metadata)
	assert.Equal(t, repoPart.CreatedAt, part.CreatedAt)
	assert.Equal(t, repoPart.UpdatedAt, part.UpdatedAt)
}

func TestModelPartInfoToRepoModel(t *testing.T) {
	info := &model.PartInfo{
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(10, 1000),
		StockQuantity: int64(gofakeit.Number(1, 100)),
		Category:      model.PartCategory(gofakeit.Number(0, 4)),
		Dimensions: &model.PartDimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(1, 100),
		},
		Manufacturer: &model.PartManufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: func() []string {
			tags := []string{}
			for i := 0; i < gofakeit.Number(1, 5); i++ {
				tags = append(tags, gofakeit.Word())
			}
			return tags
		}(),
		Metadata: func() map[string]*any {
			meta := map[string]*any{}
			for i := 0; i < gofakeit.Number(1, 5); i++ {
				value := any(gofakeit.Word())
				meta[gofakeit.Word()] = &value
			}
			return meta
		}(),
	}

	repoInfo := ModelPartInfoToRepoModel(info)

	assert.Equal(t, info.Name, repoInfo.Name)
	assert.Equal(t, info.Description, repoInfo.Description)
	assert.Equal(t, info.Price, repoInfo.Price)
	assert.Equal(t, info.StockQuantity, repoInfo.StockQuantity)
	assert.Equal(t, ModelPartCategoryToRepoModel(info.Category), repoInfo.Category)
	assert.Equal(t, info.Tags, repoInfo.Tags)
	assert.Equal(t, info.Metadata, repoInfo.Metadata)
}

func TestRepoModelPartInfoToModel(t *testing.T) {
	repoInfo := &repoModel.PartInfo{
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(10, 1000),
		StockQuantity: int64(gofakeit.Number(1, 100)),
		Category:      repoModel.PartCategory(gofakeit.Number(0, 4)),
		Dimensions: &repoModel.PartDimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(1, 100),
		},
		Manufacturer: &repoModel.PartManufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: func() []string {
			tags := []string{}
			for i := 0; i < gofakeit.Number(1, 5); i++ {
				tags = append(tags, gofakeit.Word())
			}
			return tags
		}(),
		Metadata: func() map[string]*any {
			meta := map[string]*any{}
			for i := 0; i < gofakeit.Number(1, 5); i++ {
				value := any(gofakeit.Word())
				meta[gofakeit.Word()] = &value
			}
			return meta
		}(),
	}

	info := RepoModelPartInfoToModel(repoInfo)

	assert.Equal(t, repoInfo.Name, info.Name)
	assert.Equal(t, repoInfo.Description, info.Description)
	assert.Equal(t, repoInfo.Price, info.Price)
	assert.Equal(t, repoInfo.StockQuantity, info.StockQuantity)
	assert.Equal(t, RepoModelPartCategoryToModel(repoInfo.Category), info.Category)
	assert.Equal(t, repoInfo.Tags, info.Tags)
	assert.Equal(t, repoInfo.Metadata, info.Metadata)
}

func TestModelPartCategoryToRepoModel(t *testing.T) {
	testCases := []struct {
		modelCategory     model.PartCategory
		repoModelCategory repoModel.PartCategory
	}{
		{model.PartCategory_ENGINE, repoModel.PartCategory_ENGINE},
		{model.PartCategory_FUEL, repoModel.PartCategory_FUEL},
		{model.PartCategory_PORTHOLE, repoModel.PartCategory_PORTHOLE},
		{model.PartCategory_WING, repoModel.PartCategory_WING},
		{model.PartCategory_UNSPECIFIED, repoModel.PartCategory_UNSPECIFIED},
	}

	for _, tc := range testCases {
		result := ModelPartCategoryToRepoModel(tc.modelCategory)
		assert.Equal(t, tc.repoModelCategory, result)
	}
}

func TestRepoModelPartCategoryToModel(t *testing.T) {
	testCases := []struct {
		repoModelCategory repoModel.PartCategory
		modelCategory     model.PartCategory
	}{
		{repoModel.PartCategory_ENGINE, model.PartCategory_ENGINE},
		{repoModel.PartCategory_FUEL, model.PartCategory_FUEL},
		{repoModel.PartCategory_PORTHOLE, model.PartCategory_PORTHOLE},
		{repoModel.PartCategory_WING, model.PartCategory_WING},
		{repoModel.PartCategory_UNSPECIFIED, model.PartCategory_UNSPECIFIED},
	}

	for _, tc := range testCases {
		result := RepoModelPartCategoryToModel(tc.repoModelCategory)
		assert.Equal(t, tc.modelCategory, result)
	}
}

func TestModelPartDimensionsToRepoModel(t *testing.T) {
	dimensions := &model.PartDimensions{
		Length: gofakeit.Float64Range(1, 100),
		Width:  gofakeit.Float64Range(1, 100),
		Height: gofakeit.Float64Range(1, 100),
		Weight: gofakeit.Float64Range(1, 100),
	}

	repoDimensions := ModelPartDimensionsToRepoModel(dimensions)

	assert.Equal(t, dimensions.Length, repoDimensions.Length)
	assert.Equal(t, dimensions.Width, repoDimensions.Width)
	assert.Equal(t, dimensions.Height, repoDimensions.Height)
	assert.Equal(t, dimensions.Weight, repoDimensions.Weight)
}

func TestRepoModelPartDimensionsToModel(t *testing.T) {
	repoDimensions := &repoModel.PartDimensions{
		Length: gofakeit.Float64Range(1, 100),
		Width:  gofakeit.Float64Range(1, 100),
		Height: gofakeit.Float64Range(1, 100),
		Weight: gofakeit.Float64Range(1, 100),
	}

	dimensions := RepoModelPartDimensionsToModel(repoDimensions)

	assert.Equal(t, repoDimensions.Length, dimensions.Length)
	assert.Equal(t, repoDimensions.Width, dimensions.Width)
	assert.Equal(t, repoDimensions.Height, dimensions.Height)
	assert.Equal(t, repoDimensions.Weight, dimensions.Weight)
}

func TestModelPartManufacturerToRepoModel(t *testing.T) {
	manufacturer := &model.PartManufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}

	repoManufacturer := ModelPartManufacturerToRepoModel(manufacturer)

	assert.Equal(t, manufacturer.Name, repoManufacturer.Name)
	assert.Equal(t, manufacturer.Country, repoManufacturer.Country)
	assert.Equal(t, manufacturer.Website, repoManufacturer.Website)
}

func TestRepoModelPartManufacturerToModel(t *testing.T) {
	repoManufacturer := &repoModel.PartManufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}

	manufacturer := RepoModelPartManufacturerToModel(repoManufacturer)

	assert.Equal(t, repoManufacturer.Name, manufacturer.Name)
	assert.Equal(t, repoManufacturer.Country, manufacturer.Country)
	assert.Equal(t, repoManufacturer.Website, manufacturer.Website)
}
