package kafka

import (
	"context"

	"github.com/Reensef/go-microservices-course/assembly/internal/model"
)

type OrderConsumer interface {
	RunConsumer(ctx context.Context) error
}

type ShipProducer interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
