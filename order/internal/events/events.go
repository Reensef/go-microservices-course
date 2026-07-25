package kafka

import (
	"context"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

type ShipConsumer interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducer interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
