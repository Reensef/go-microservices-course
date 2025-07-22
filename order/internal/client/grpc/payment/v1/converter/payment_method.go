package converter

import (
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func ToProtoPaymentMethod(paymentMethod model.OrderPaymentMethod) paymentV1.PaymentMethod {
	switch paymentMethod {
	case model.OrderPaymentMethod_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.OrderPaymentMethod_INVESTOR_MONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.OrderPaymentMethod_CREDIT_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.OrderPaymentMethod_SBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
