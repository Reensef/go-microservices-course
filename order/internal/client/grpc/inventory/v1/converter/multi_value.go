package converter

import (
	"github.com/Reensef/go-microservices-course/platform/pkg/multivalue"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func ToProtoMultiValue(multiValue *multivalue.MultiValue) *inventoryGrpc.Value {
	switch multiValue.Type() {
	case multivalue.MultiValueType_INT64:
		return &inventoryGrpc.Value{
			Kind: &inventoryGrpc.Value_Int64Value{Int64Value: multiValue.MustInt64()},
		}
	case multivalue.MultiValueType_FLOAT64:
		return &inventoryGrpc.Value{
			Kind: &inventoryGrpc.Value_DoubleValue{DoubleValue: multiValue.MustFloat64()},
		}
	case multivalue.MultiValueType_BOOL:
		return &inventoryGrpc.Value{
			Kind: &inventoryGrpc.Value_BoolValue{BoolValue: multiValue.MustBool()},
		}
	case multivalue.MultiValueType_STRING:
		return &inventoryGrpc.Value{
			Kind: &inventoryGrpc.Value_StringValue{StringValue: multiValue.MustString()},
		}
	default:
		return &inventoryGrpc.Value{}
	}
}

func ToMultiValue(value *inventoryGrpc.Value) *multivalue.MultiValue {
	switch value.GetKind().(type) {
	case *inventoryGrpc.Value_Int64Value:
		multiValue := multivalue.MultiValue{}
		multiValue.SetInt64(value.GetInt64Value())
		return &multiValue
	case *inventoryGrpc.Value_DoubleValue:
		multiValue := multivalue.MultiValue{}
		multiValue.SetFloat64(value.GetDoubleValue())
		return &multiValue
	case *inventoryGrpc.Value_BoolValue:
		multiValue := multivalue.MultiValue{}
		multiValue.SetBool(value.GetBoolValue())
		return &multiValue
	case *inventoryGrpc.Value_StringValue:
		multiValue := multivalue.MultiValue{}
		multiValue.SetString(value.GetStringValue())
		return &multiValue
	default:
		return nil
	}
}
