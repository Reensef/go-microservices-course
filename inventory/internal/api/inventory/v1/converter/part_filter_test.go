package converter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func TestToModelPartsFilter(t *testing.T) {
	protoFilter := &inventoryV1.PartsFilter{
		Uuids: []string{
			"f2bd04cb-a838-4de0-ade2-cdbb2b389167",
			"ea815eb8-a569-47fb-89fb-3052b3612396",
		},
		Names: []string{
			"name1",
			"name2",
		},
		Categories: []inventoryV1.Category{
			inventoryV1.Category_CATEGORY_ENGINE,
			inventoryV1.Category_CATEGORY_FUEL,
		},
		ManufacturerCountries: []string{
			"country1",
			"country2",
		},
		Tags: []string{
			"tag1",
			"tag2",
		},
	}

	modelFilter, err := ToModelPartsFilter(protoFilter)

	assert.NoError(t, err)
	assert.Equal(t, modelFilter.Uuids, []uuid.UUID{
		uuid.MustParse("f2bd04cb-a838-4de0-ade2-cdbb2b389167"),
		uuid.MustParse("ea815eb8-a569-47fb-89fb-3052b3612396"),
	})
	assert.Equal(t, modelFilter.Names, []string{"name1", "name2"})
	assert.Equal(t, modelFilter.Categories, []model.PartCategory{
		model.PartCategory_ENGINE,
		model.PartCategory_FUEL,
	})
	assert.Equal(t, modelFilter.ManufacturerCountries, []string{"country1", "country2"})
	assert.Equal(t, modelFilter.Tags, []string{"tag1", "tag2"})
}

func TestToModelPartsFilter_Nil(t *testing.T) {
	modelFilter, err := ToModelPartsFilter(nil)

	assert.Nil(t, modelFilter)
	assert.NoError(t, err)
}
