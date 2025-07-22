package converter

import (
	model "github.com/Reensef/go-microservices-course/payment/internal/model"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func ToModelPaymentMethod(paymentMethod paymentV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethod_CARD
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethod_INVESTOR_MONEY
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethod_CREDIT_CARD
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethod_SBP
	default:
		return model.PaymentMethod_UNSPECIFIED
	}
}

func ToProtoPaymentMethod(paymentMethod model.PaymentMethod) paymentV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethod_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethod_INVESTOR_MONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.PaymentMethod_CREDIT_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethod_SBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
