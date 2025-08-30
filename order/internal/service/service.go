package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, info *model.OrderInfo) (*model.Order, error)
	GetOrderByUUID(ctx context.Context, orderUuid string) (*model.Order, error)
	CancelOrder(ctx context.Context, orderUuid string) error
	PayOrder(
		ctx context.Context,
		orderUuid string,
		userUuid string,
		paymentMethod model.OrderPaymentMethod,
	) (transactionUUID *string, err error)
}
