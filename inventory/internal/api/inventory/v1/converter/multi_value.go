package converter

import (
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

func ToProtoMultiValue(multiValue *utils.MultiValue) *inventoryV1.Value {
	switch multiValue.Type() {
	case utils.MultiValueType_INT64:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{Int64Value: multiValue.MustInt64()},
		}
	case utils.MultiValueType_FLOAT64:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{DoubleValue: multiValue.MustFloat64()},
		}
	case utils.MultiValueType_BOOL:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{BoolValue: multiValue.MustBool()},
		}
	case utils.MultiValueType_STRING:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{StringValue: multiValue.MustString()},
		}
	default:
		return &inventoryV1.Value{}
	}
}

func ToMultiValue(value *inventoryV1.Value) *utils.MultiValue {
	switch value.GetKind().(type) {
	case *inventoryV1.Value_Int64Value:
		multiValue := utils.MultiValue{}
		multiValue.SetInt64(value.GetInt64Value())
		return &multiValue
	case *inventoryV1.Value_DoubleValue:
		multiValue := utils.MultiValue{}
		multiValue.SetFloat64(value.GetDoubleValue())
		return &multiValue
	case *inventoryV1.Value_BoolValue:
		multiValue := utils.MultiValue{}
		multiValue.SetBool(value.GetBoolValue())
		return &multiValue
	case *inventoryV1.Value_StringValue:
		multiValue := utils.MultiValue{}
		multiValue.SetString(value.GetStringValue())
		return &multiValue
	default:
		return nil
	}
}
