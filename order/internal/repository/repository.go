package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(
		ctx context.Context,
		info *model.OrderInfo,
	) (*model.Order, error)

	GetOrderByUUID(
		ctx context.Context,
		orderUuid uuid.UUID,
	) (*model.Order, error)

	CancelOrder(
		ctx context.Context,
		orderUuid uuid.UUID,
	) error

	PayOrder(
		ctx context.Context,
		orderUuid uuid.UUID,
		transactionUUID uuid.UUID,
		paymentMethod model.OrderPaymentMethod,
	) error
}
