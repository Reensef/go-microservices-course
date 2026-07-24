package assembly

import (
	events "github.com/Reensef/go-microservices-course/assembly/internal/events"
)

type service struct {
	shipProducer events.ShipProducer
}

func NewService(shipProducer events.ShipProducer) *service {
	return &service{
		shipProducer: shipProducer,
	}
}
