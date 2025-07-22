package converter

import (
	"testing"

	model "github.com/Reensef/go-microservices-course/payment/internal/model"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func TestToModelPaymentMethod(t *testing.T) {
	tests := []struct {
		name                string
		paymentMethod       paymentV1.PaymentMethod
		expectedModelMethod model.PaymentMethod
	}{
		{
			name:                "Test CARD payment method",
			paymentMethod:       paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			expectedModelMethod: model.PaymentMethod_CARD,
		},
		{
			name:                "Test INVESTOR_MONEY payment method",
			paymentMethod:       paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
			expectedModelMethod: model.PaymentMethod_INVESTOR_MONEY,
		},
		{
			name:                "Test CREDIT_CARD payment method",
			paymentMethod:       paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
			expectedModelMethod: model.PaymentMethod_CREDIT_CARD,
		},
		{
			name:                "Test SBP payment method",
			paymentMethod:       paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
			expectedModelMethod: model.PaymentMethod_SBP,
		},
		{
			name:                "Test unknown payment method",
			paymentMethod:       paymentV1.PaymentMethod(100),
			expectedModelMethod: model.PaymentMethod_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelMethod := ToModelPaymentMethod(tt.paymentMethod)
			if modelMethod != tt.expectedModelMethod {
				t.Errorf("ToModelPaymentMethod() = %v, want %v", modelMethod, tt.expectedModelMethod)
			}
		})
	}
}

func TestToProtoPaymentMethod(t *testing.T) {
	tests := []struct {
		name                string
		paymentMethod       model.PaymentMethod
		expectedProtoMethod paymentV1.PaymentMethod
	}{
		{
			name:                "Test CARD payment method",
			paymentMethod:       model.PaymentMethod_CARD,
			expectedProtoMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		},
		{
			name:                "Test INVESTOR_MONEY payment method",
			paymentMethod:       model.PaymentMethod_INVESTOR_MONEY,
			expectedProtoMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
		},
		{
			name:                "Test CREDIT_CARD payment method",
			paymentMethod:       model.PaymentMethod_CREDIT_CARD,
			expectedProtoMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		},
		{
			name:                "Test SBP payment method",
			paymentMethod:       model.PaymentMethod_SBP,
			expectedProtoMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		},
		{
			name:                "Test unknown payment method",
			paymentMethod:       model.PaymentMethod(100),
			expectedProtoMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoMethod := ToProtoPaymentMethod(tt.paymentMethod)
			if protoMethod != tt.expectedProtoMethod {
				t.Errorf("ToProtoPaymentMethod() = %v, want %v", protoMethod, tt.expectedProtoMethod)
			}
		})
	}
}
