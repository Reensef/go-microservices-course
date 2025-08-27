package converter

import (
	"github.com/Reensef/go-microservices-course/platform/pkg/multivalue"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToProtoMultiValue(multiValue *multivalue.MultiValue) *inventoryV1.Value {
	switch multiValue.Type() {
	case multivalue.MultiValueType_INT64:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{Int64Value: multiValue.MustInt64()},
		}
	case multivalue.MultiValueType_FLOAT64:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{DoubleValue: multiValue.MustFloat64()},
		}
	case multivalue.MultiValueType_BOOL:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{BoolValue: multiValue.MustBool()},
		}
	case multivalue.MultiValueType_STRING:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{StringValue: multiValue.MustString()},
		}
	default:
		return &inventoryV1.Value{}
	}
}

func ToMultiValue(value *inventoryV1.Value) *multivalue.MultiValue {
	switch value.GetKind().(type) {
	case *inventoryV1.Value_Int64Value:
		multiValue := multivalue.MultiValue{}
		multiValue.SetInt64(value.GetInt64Value())
		return &multiValue
	case *inventoryV1.Value_DoubleValue:
		multiValue := multivalue.MultiValue{}
		multiValue.SetFloat64(value.GetDoubleValue())
		return &multiValue
	case *inventoryV1.Value_BoolValue:
		multiValue := multivalue.MultiValue{}
		multiValue.SetBool(value.GetBoolValue())
		return &multiValue
	case *inventoryV1.Value_StringValue:
		multiValue := multivalue.MultiValue{}
		multiValue.SetString(value.GetStringValue())
		return &multiValue
	default:
		return nil
	}
}
