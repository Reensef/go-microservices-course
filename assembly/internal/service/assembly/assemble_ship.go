package assembly

import (
	"context"
	"time"

	"github.com/Reensef/go-microservices-course/assembly/internal/model"
	"github.com/google/uuid"
)

func (s *service) AssembleShip(
	ctx context.Context, orderUuid string, userUuid string,
) error {
	timeNow := time.Now()
	// emulation of ship assembly
	time.Sleep(5 * time.Second)
	buildDuration := time.Since(timeNow)

	err := s.shipProducer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		UUID:          uuid.New().String(),
		OrderUUID:     orderUuid,
		UserUUID:      userUuid,
		BuildDuration: buildDuration,
	})
	if err != nil {

	}
	return nil
}
