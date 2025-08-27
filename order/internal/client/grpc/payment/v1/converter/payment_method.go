package converter

import (
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	paymentGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

func ToProtoPaymentMethod(paymentMethod model.OrderPaymentMethod) paymentGrpc.PaymentMethod {
	switch paymentMethod {
	case model.OrderPaymentMethod_CARD:
		return paymentGrpc.PaymentMethod_PAYMENT_METHOD_CARD
	case model.OrderPaymentMethod_INVESTOR_MONEY:
		return paymentGrpc.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.OrderPaymentMethod_CREDIT_CARD:
		return paymentGrpc.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.OrderPaymentMethod_SBP:
		return paymentGrpc.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return paymentGrpc.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
