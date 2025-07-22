package converter

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

func TestToRepoModelPart(t *testing.T) {
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
				for range 5 {
					tags = append(tags, gofakeit.Word())
				}
				return tags
			}(),
			Metadata: map[string]utils.MultiValue{
				"key1": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetString("value")
					return value
				}(),
				"key2": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetInt64(10)
					return value
				}(),
				"key3": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetFloat64(10.10)
					return value
				}(),
				"key4": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetBool(true)
					return value
				}(),
			},
		},
	}

	repoPart := ToRepoModelPart(part)

	assert.Equal(t, part.Uuid, repoPart.Uuid)
	assert.Equal(t, part.Info.Name, repoPart.Info.Name)
	assert.Equal(t, part.Info.Description, repoPart.Info.Description)
	assert.Equal(t, part.Info.Price, repoPart.Info.Price)
	assert.Equal(t, part.Info.StockQuantity, repoPart.Info.StockQuantity)
	assert.Equal(t, ToRepoModelPartCategory(part.Info.Category), repoPart.Info.Category)
	assert.Equal(t, part.Info.Tags, repoPart.Info.Tags)
	assert.Equal(t, part.Info.Metadata, repoPart.Info.Metadata)
	assert.Equal(t, part.CreatedAt, repoPart.CreatedAt)
	assert.Equal(t, part.UpdatedAt, repoPart.UpdatedAt)
}

func TestToModelPart(t *testing.T) {
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
			Metadata: map[string]utils.MultiValue{
				"key1": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetString("value")
					return value
				}(),
				"key2": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetInt64(10)
					return value
				}(),
				"key3": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetFloat64(10.10)
					return value
				}(),
				"key4": func() utils.MultiValue {
					value := utils.MultiValue{}
					value.SetBool(true)
					return value
				}(),
			},
		},
	}

	part := ToModelPart(repoPart)

	assert.Equal(t, repoPart.Uuid, part.Uuid)
	assert.Equal(t, repoPart.Info.Name, part.Info.Name)
	assert.Equal(t, repoPart.Info.Description, part.Info.Description)
	assert.Equal(t, repoPart.Info.Price, part.Info.Price)
	assert.Equal(t, repoPart.Info.StockQuantity, part.Info.StockQuantity)
	assert.Equal(t, ToModelPartCategory(repoPart.Info.Category), part.Info.Category)
	assert.Equal(t, repoPart.Info.Tags, part.Info.Tags)
	assert.Equal(t, repoPart.Info.Metadata, part.Info.Metadata)
	assert.Equal(t, repoPart.CreatedAt, part.CreatedAt)
	assert.Equal(t, repoPart.UpdatedAt, part.UpdatedAt)
}

func TestToRepoModelPartCategory(t *testing.T) {
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
		result := ToRepoModelPartCategory(tc.modelCategory)
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
		result := ToModelPartCategory(tc.repoModelCategory)
		assert.Equal(t, tc.modelCategory, result)
	}
}
