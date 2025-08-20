package repository

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(
		ctx context.Context,
		info *model.OrderInfo,
	) (*model.Order, error)

	GetOrderByUUID(
		ctx context.Context,
		orderUuid string,
	) (*model.Order, error)

	GetOrderStatus(
		ctx context.Context,
		orderUuid string,
	) (model.OrderStatus, error)

	CancelOrder(
		ctx context.Context,
		orderUuid string,
	) error

	PayOrder(
		ctx context.Context,
		orderUuid string,
		transactionUUID string,
		paymentMethod model.OrderPaymentMethod,
	) error
}
