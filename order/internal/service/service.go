package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUuid uuid.UUID, partUuids []uuid.UUID) (model.Order, error)
	GetOrderByUUID(ctx context.Context, orderUuid uuid.UUID) (model.Order, error)
	CancelOrder(ctx context.Context, orderUuid uuid.UUID) error
	PayOrder(ctx context.Context, orderUuid uuid.UUID, paymentMethod model.OrderPaymentMethod) (uuid.UUID, error)
}
