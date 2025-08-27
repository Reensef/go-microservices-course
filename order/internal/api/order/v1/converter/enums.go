package converter

import (
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func ToModelPaymentMethod(paymentMethod orderApi.PaymentMethod) model.OrderPaymentMethod {
	switch paymentMethod {
	case orderApi.PaymentMethodCARD:
		return model.OrderPaymentMethod_CARD
	case orderApi.PaymentMethodCREDITCARD:
		return model.OrderPaymentMethod_CREDIT_CARD
	case orderApi.PaymentMethodSBP:
		return model.OrderPaymentMethod_SBP
	case orderApi.PaymentMethodINVESTORMONEY:
		return model.OrderPaymentMethod_INVESTOR_MONEY
	default:
		return model.OrderPaymentMethod_UNSPECIFIED
	}
}

func ToApiPaymentMethod(paymentMethod model.OrderPaymentMethod) orderApi.PaymentMethod {
	switch paymentMethod {
	case model.OrderPaymentMethod_CARD:
		return orderApi.PaymentMethodCARD
	case model.OrderPaymentMethod_CREDIT_CARD:
		return orderApi.PaymentMethodCREDITCARD
	case model.OrderPaymentMethod_SBP:
		return orderApi.PaymentMethodSBP
	case model.OrderPaymentMethod_INVESTOR_MONEY:
		return orderApi.PaymentMethodINVESTORMONEY
	default:
		return orderApi.PaymentMethodUNSPECIFIED
	}
}

func ToModelOrderStatus(status orderApi.OrderStatus) model.OrderStatus {
	switch status {
	case orderApi.OrderStatusCANCELED:
		return model.OrderStatus_CANCELED
	case orderApi.OrderStatusPAID:
		return model.OrderStatus_PAID
	default:
		return model.OrderStatus_PENDING_PAYMENT
	}
}

func ToApiOrderStatus(status model.OrderStatus) orderApi.OrderStatus {
	switch status {
	case model.OrderStatus_PAID:
		return orderApi.OrderStatusPAID
	case model.OrderStatus_CANCELED:
		return orderApi.OrderStatusCANCELED
	default:
		return orderApi.OrderStatusPENDINGPAYMENT
	}
}
