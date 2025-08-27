package converter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/multivalue"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func TesToProtoPart(t *testing.T) {
	part := &model.Part{
		ID:        uuid.New().String(),
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
			Metadata: map[string]multivalue.MultiValue{
				"key1": func() multivalue.MultiValue {
					value := multivalue.MultiValue{}
					value.SetString("value")
					return value
				}(),
				"key2": func() multivalue.MultiValue {
					value := multivalue.MultiValue{}
					value.SetInt64(10)
					return value
				}(),
				"key3": func() multivalue.MultiValue {
					value := multivalue.MultiValue{}
					value.SetFloat64(10.10)
					return value
				}(),
				"key4": func() multivalue.MultiValue {
					value := multivalue.MultiValue{}
					value.SetBool(true)
					return value
				}(),
			},
		},
	}

	protoPart := ToProtoPart(part)

	assert.Equal(t, part.ID, protoPart.Id)
	assert.Equal(t, part.Info.Name, protoPart.Name)
	assert.Equal(t, part.Info.Description, protoPart.Description)
	assert.Equal(t, part.Info.Price, protoPart.Price)
	assert.Equal(t, part.Info.StockQuantity, protoPart.StockQuantity)
	assert.Equal(t, ToProtoPartCategory(part.Info.Category), protoPart.Category)
	assert.Equal(t, part.Info.Tags, protoPart.Tags)

	for key, protoMeta := range protoPart.Metadata {
		meta, ok := part.Info.Metadata[key]
		if !ok {
			assert.Fail(t, "metadata key not found")
			continue
		}

		assert.Equal(t, ToProtoMultiValue(&meta), protoMeta)
	}

	assert.Equal(t, part.CreatedAt.UTC(), protoPart.CreatedAt.AsTime())
	assert.Equal(t, part.UpdatedAt.UTC(), protoPart.UpdatedAt.AsTime())
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
			"ket2": func() *inventoryV1.Value {
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

	modelPart := ToModelPart(protoPart)

	assert.Equal(t, modelPart.ID, protoPart.Id)
	assert.Equal(t, modelPart.Info.Name, protoPart.Name)
	assert.Equal(t, modelPart.Info.Description, protoPart.Description)
	assert.Equal(t, modelPart.Info.Price, protoPart.Price)
	assert.Equal(t, modelPart.Info.StockQuantity, protoPart.StockQuantity)
	assert.Equal(t, modelPart.Info.Category, ToModelPartCategory(protoPart.Category))
	assert.Equal(t, modelPart.Info.Tags, protoPart.Tags)

	for key, protoMeta := range protoPart.Metadata {
		modelMeta, ok := modelPart.Info.Metadata[key]
		if !ok {
			assert.Fail(t, "metadata key not found")
			continue
		}

		assert.Equal(t, &modelMeta, ToMultiValue(protoMeta))
	}

	assert.Equal(t, modelPart.CreatedAt.UTC(), protoPart.CreatedAt.AsTime())
	assert.Equal(t, modelPart.UpdatedAt.UTC(), protoPart.UpdatedAt.AsTime())
}
