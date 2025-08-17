package converter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/Reensef/go-microservices-course/order/internal/model"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

func TestToProtoPart(t *testing.T) {
	part := &model.Part{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Info: model.PartInfo{
			Name:          "Name",
			Description:   "Description",
			Price:         10.0,
			StockQuantity: 50,
			Category:      model.PartCategory_ENGINE,
			Dimensions: &model.PartDimensions{
				Length: 1.0,
				Width:  2.0,
				Height: 3.0,
				Weight: 4.0,
			},
			Manufacturer: &model.PartManufacturer{
				Name:    "MName",
				Country: "Country",
				Website: "Website",
			},
			Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
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

	expected := &inventoryV1.Part{
		Name:          part.Info.Name,
		Id:            part.Id,
		Description:   part.Info.Description,
		Price:         part.Info.Price,
		StockQuantity: part.Info.StockQuantity,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: part.Info.Dimensions.Length,
			Width:  part.Info.Dimensions.Width,
			Height: part.Info.Dimensions.Height,
			Weight: part.Info.Dimensions.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    part.Info.Manufacturer.Name,
			Country: part.Info.Manufacturer.Country,
			Website: part.Info.Manufacturer.Website,
		},
		Tags: part.Info.Tags,
		Metadata: map[string]*inventoryV1.Value{
			"key1": {
				Kind: &inventoryV1.Value_StringValue{
					StringValue: "value",
				},
			},
			"key2": {
				Kind: &inventoryV1.Value_Int64Value{
					Int64Value: 10,
				},
			},
			"key3": {
				Kind: &inventoryV1.Value_DoubleValue{
					DoubleValue: 10.10,
				},
			},
			"key4": {
				Kind: &inventoryV1.Value_BoolValue{
					BoolValue: true,
				},
			},
		},
		CreatedAt: timestamppb.New(part.CreatedAt),
		UpdatedAt: timestamppb.New(part.UpdatedAt),
	}

	protoPart := ToProtoPart(part)
	require.NotNil(t, protoPart)

	assert.Equal(t, expected, protoPart)
}

func TestToModelPart(t *testing.T) {
	protoPart := &inventoryV1.Part{
		Id:        uuid.New().String(),
		CreatedAt: timestamppb.New(time.Now().Add(10)),
		UpdatedAt: timestamppb.New(time.Now().Add(20)),

		Name:          "Name",
		Description:   "Description",
		Price:         10.0,
		StockQuantity: 50,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 1.0,
			Width:  2.0,
			Height: 3.0,
			Weight: 4.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "MName",
			Country: "Country",
			Website: "Website",
		},
		Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		Metadata: map[string]*inventoryV1.Value{
			"key1": func() *inventoryV1.Value {
				value := &inventoryV1.Value{}
				value.Kind = &inventoryV1.Value_StringValue{StringValue: "value"}
				return value
			}(),
			"key2": func() *inventoryV1.Value {
				value := &inventoryV1.Value{}
				value.Kind = &inventoryV1.Value_Int64Value{Int64Value: 10}
				return value
			}(),
			"key3": func() *inventoryV1.Value {
				value := &inventoryV1.Value{}
				value.Kind = &inventoryV1.Value_DoubleValue{DoubleValue: 10.10}
				return value
			}(),
			"key4": func() *inventoryV1.Value {
				value := &inventoryV1.Value{}
				value.Kind = &inventoryV1.Value_BoolValue{BoolValue: true}
				return value
			}(),
		},
	}

	expected := &model.Part{
		Id:        protoPart.GetId(),
		CreatedAt: protoPart.CreatedAt.AsTime().UTC(),
		UpdatedAt: protoPart.UpdatedAt.AsTime().UTC(),
		Info: model.PartInfo{
			Name:          protoPart.Name,
			Description:   protoPart.Description,
			Price:         protoPart.Price,
			StockQuantity: protoPart.StockQuantity,
			Category:      model.PartCategory_ENGINE,
			Tags:          protoPart.Tags,
			Dimensions: &model.PartDimensions{
				Length: protoPart.Dimensions.Length,
				Width:  protoPart.Dimensions.Width,
				Height: protoPart.Dimensions.Height,
				Weight: protoPart.Dimensions.Weight,
			},
			Manufacturer: &model.PartManufacturer{
				Name:    protoPart.Manufacturer.Name,
				Country: protoPart.Manufacturer.Country,
				Website: protoPart.Manufacturer.Website,
			},
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

	modelPart, err := ToModelPart(protoPart)
	require.Nil(t, err)
	assert.Equal(t, expected, modelPart)
}
