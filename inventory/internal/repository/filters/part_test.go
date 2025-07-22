package filters

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func TestMatchPartFilters(t *testing.T) {
	tagsCount := 5
	part := &repoModel.Part{
		Info: repoModel.PartInfo{
			Name:     gofakeit.ProductName(),
			Category: repoModel.PartCategory_ENGINE,
			Manufacturer: &repoModel.PartManufacturer{
				Country: gofakeit.Country(),
			},
			Tags: func() []string {
				tags := []string{}
				for range tagsCount {
					tags = append(tags, gofakeit.Word())
				}
				return tags
			}(),
		},
	}

	filter := &model.PartsFilter{
		Names:                 []string{part.Info.Name},
		Categories:            []model.PartCategory{model.PartCategory_ENGINE},
		ManufacturerCountries: []string{part.Info.Manufacturer.Country},
		Tags:                  []string{part.Info.Tags[gofakeit.IntN(tagsCount)]},
	}

	assert.True(t, MatchPartFilters(part, filter))

	filter.Names = []string{"nonexistent"}
	assert.False(t, MatchPartFilters(part, filter))
}

func TestFilterPartByNames(t *testing.T) {
	part := &repoModel.Part{
		Info: repoModel.PartInfo{
			Name: gofakeit.ProductName(),
		},
	}

	// Пустой фильтр (всегда true)
	assert.True(t, filterPartByNames(part, []string{}))

	// Подходящее имя
	assert.True(t, filterPartByNames(part, []string{part.Info.Name}))

	// Неподходящее имя
	assert.False(t, filterPartByNames(part, []string{"nonexistent"}))
}

func TestFilterPartByCategories(t *testing.T) {
	part := &repoModel.Part{
		Info: repoModel.PartInfo{
			Category: repoModel.PartCategory_ENGINE,
		},
	}

	// Пустой фильтр (всегда true)
	assert.True(t, filterPartByCategories(part, []model.PartCategory{}))

	// Подходящая категория
	assert.True(t, filterPartByCategories(part, []model.PartCategory{model.PartCategory_ENGINE}))

	// Неподходящая категория
	assert.False(t, filterPartByCategories(part, []model.PartCategory{model.PartCategory_FUEL}))
}

func TestFilterPartByManufacturerCountries(t *testing.T) {
	part := &repoModel.Part{
		Info: repoModel.PartInfo{
			Manufacturer: &repoModel.PartManufacturer{
				Country: "USA",
			},
		},
	}

	// Пустой фильтр (всегда true)
	assert.True(t, filterPartByManufacturerCountries(part, []string{}))

	// Подходящая страна
	assert.True(t, filterPartByManufacturerCountries(part, []string{part.Info.Manufacturer.Country}))

	// Неподходящая страна
	assert.False(t, filterPartByManufacturerCountries(part, []string{"Canada"}))
}

func TestFilterPartByTags(t *testing.T) {
	tagsCount := 5
	part := &repoModel.Part{
		Info: repoModel.PartInfo{
			Tags: func() []string {
				tags := []string{}
				for range tagsCount {
					tags = append(tags, gofakeit.Word())
				}
				return tags
			}(),
		},
	}

	// Пустой фильтр (всегда true)
	assert.True(t, filterPartByTags(part, []string{}))

	// Подходящий тег
	assert.True(t, filterPartByTags(part, []string{part.Info.Tags[gofakeit.IntN(tagsCount)]}))

	// Неподходящий тег
	assert.False(t, filterPartByTags(part, []string{"noexistent"}))

	// Множественные теги (хотя бы один совпадает)
	assert.True(t, filterPartByTags(part, []string{
		part.Info.Tags[gofakeit.IntN(tagsCount)],
		part.Info.Tags[gofakeit.IntN(tagsCount)],
	}),
	)
}
