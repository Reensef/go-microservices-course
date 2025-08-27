package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func TestToProtoFilter(t *testing.T) {
	filter := model.PartsFilter{
		Ids: []string{
			"3f74bbbe-ad48-43cd-bffd-80d0a1b6a8e6",
			"b8a251c6-df82-4add-bce0-199a81d39cdf",
		},
		Names: []string{
			"name1",
			"name2",
		},
		Categories: []model.PartCategory{
			model.PartCategory_ENGINE,
			model.PartCategory_FUEL,
			model.PartCategory_PORTHOLE,
			model.PartCategory_WING,
			model.PartCategory_UNSPECIFIED,
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

	expected := &inventoryGrpc.PartsFilter{
		Ids: []string{
			"3f74bbbe-ad48-43cd-bffd-80d0a1b6a8e6",
			"b8a251c6-df82-4add-bce0-199a81d39cdf",
		},
		Names: filter.Names,
		Categories: []inventoryGrpc.Category{
			inventoryGrpc.Category_CATEGORY_ENGINE,
			inventoryGrpc.Category_CATEGORY_FUEL,
			inventoryGrpc.Category_CATEGORY_PORTHOLE,
			inventoryGrpc.Category_CATEGORY_WING,
			inventoryGrpc.Category_CATEGORY_UNSPECIFIED,
		},
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}

	protoFilter := ToProtoFilter(&filter)
	require.NotNil(t, protoFilter)
	assert.Equal(t, expected, protoFilter)
}
