package converter

import (
	"github.com/Reensef/go-microservices-course/order/internal/model"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

func APIPaymentMethodToModel(paymentMethod orderV1.PaymentMethod) model.OrderPaymentMethod {
	switch paymentMethod {
	case orderV1.PaymentMethodCARD:
		return model.OrderPaymentMethod_CARD
	case orderV1.PaymentMethodCREDITCARD:
		return model.OrderPaymentMethod_CREDIT_CARD
	case orderV1.PaymentMethodSBP:
		return model.OrderPaymentMethod_SBP
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.OrderPaymentMethod_INVESTOR_MONEY
	default:
		return model.OrderPaymentMethod_UNKNOWN
	}
}

func ModelPaymentMethodToAPI(paymentMethod model.OrderPaymentMethod) orderV1.PaymentMethod {
	switch paymentMethod {
	case model.OrderPaymentMethod_CARD:
		return orderV1.PaymentMethodCARD
	case model.OrderPaymentMethod_CREDIT_CARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.OrderPaymentMethod_SBP:
		return orderV1.PaymentMethodSBP
	case model.OrderPaymentMethod_INVESTOR_MONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func APIOrderStatusToModel(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusCANCELED:
		return model.OrderStatus_CANCELED
	case orderV1.OrderStatusPAID:
		return model.OrderStatus_PAID
	case orderV1.OrderStatusPENDINGPAYMENT:
		return model.OrderStatus_PENDING_PAYMENT
	default:
		return model.OrderPaymentMethod_UNKNOWN
	}
}

func ModelPaymentMethodToAPI(paymentMethod model.OrderPaymentMethod) orderV1.PaymentMethod {
	switch paymentMethod {
	case model.OrderPaymentMethod_CARD:
		return orderV1.PaymentMethodCARD
	case model.OrderPaymentMethod_CREDIT_CARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.OrderPaymentMethod_SBP:
		return orderV1.PaymentMethodSBP
	case model.OrderPaymentMethod_INVESTOR_MONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
