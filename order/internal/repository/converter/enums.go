package converter

import (
	model "github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func ToRepoModelPaymentMethod(paymentMethod model.OrderPaymentMethod) repoModel.OrderPaymentMethod {
	switch paymentMethod {
	case model.OrderPaymentMethod_CARD:
		return repoModel.OrderPaymentMethod_CARD
	case model.OrderPaymentMethod_CREDIT_CARD:
		return repoModel.OrderPaymentMethod_CREDIT_CARD
	case model.OrderPaymentMethod_SBP:
		return repoModel.OrderPaymentMethod_SBP
	case model.OrderPaymentMethod_INVESTOR_MONEY:
		return repoModel.OrderPaymentMethod_INVESTOR_MONEY
	default:
		return repoModel.OrderPaymentMethod_UNSPECIFIED
	}
}

func ToModelPaymentMethod(paymentMethod repoModel.OrderPaymentMethod) model.OrderPaymentMethod {
	switch paymentMethod {
	case repoModel.OrderPaymentMethod_CARD:
		return model.OrderPaymentMethod_CARD
	case repoModel.OrderPaymentMethod_CREDIT_CARD:
		return model.OrderPaymentMethod_CREDIT_CARD
	case repoModel.OrderPaymentMethod_SBP:
		return model.OrderPaymentMethod_SBP
	case repoModel.OrderPaymentMethod_INVESTOR_MONEY:
		return model.OrderPaymentMethod_INVESTOR_MONEY
	default:
		return model.OrderPaymentMethod_UNSPECIFIED
	}
}

func ToModelOrderStatus(status repoModel.OrderStatus) model.OrderStatus {
	switch status {
	case repoModel.OrderStatus_CANCELED:
		return model.OrderStatus_CANCELED
	case repoModel.OrderStatus_PAID:
		return model.OrderStatus_PAID
	case repoModel.OrderStatus_PENDING_PAYMENT:
		return model.OrderStatus_PENDING_PAYMENT
	default:
		return model.OrderStatus_UNSPECIFIED
	}
}

func ToRepoModelOrderStatus(status model.OrderStatus) repoModel.OrderStatus {
	switch status {
	case model.OrderStatus_PAID:
		return repoModel.OrderStatus_PAID
	case model.OrderStatus_CANCELED:
		return repoModel.OrderStatus_CANCELED
	case model.OrderStatus_PENDING_PAYMENT:
		return repoModel.OrderStatus_PENDING_PAYMENT
	default:
		return repoModel.OrderStatus_UNSPECIFIED
	}
}
