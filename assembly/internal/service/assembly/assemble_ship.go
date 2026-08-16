package assembly

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Reensef/go-microservices-course/assembly/internal/metric"
	"github.com/Reensef/go-microservices-course/assembly/internal/model"
)

func (s *service) AssembleShip(
	ctx context.Context, orderUuid, userUuid string,
) error {
	timeNow := time.Now()

	// emulation of ship assembly
	wait := time.After(5 * time.Second)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wait:
	}

	buildDuration := time.Since(timeNow)

	err := s.shipProducer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		UUID:          uuid.New().String(),
		OrderUUID:     orderUuid,
		UserUUID:      userUuid,
		BuildDuration: buildDuration,
	})
	if err != nil {
		metric.IncOrdersProcessed(ctx, "error")
		metric.ObserveAssemblyDuration(ctx, "error", buildDuration.Seconds())
		return err
	}

	metric.IncOrdersProcessed(ctx, "success")
	metric.ObserveAssemblyDuration(ctx, "success", buildDuration.Seconds())
	return nil
}
