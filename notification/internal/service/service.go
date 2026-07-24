package service

import (
	"context"

	"github.com/Reensef/go-microservices-course/notification/internal/model"
)

type NotificationService interface {
	NotifyOrderPaid(ctx context.Context, event model.OrderPaidEvent) error

	NotifyShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
