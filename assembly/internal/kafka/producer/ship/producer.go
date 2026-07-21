package ship

import (
	"context"

	"github.com/Reensef/go-microservices-course/assembly/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/kafka"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

type producer struct {
	kafkaProducer kafka.Producer
}

func NewProducer(kafkaProducer kafka.Producer) *producer {
	return &producer{
		kafkaProducer: kafkaProducer,
	}
}

func (p *producer) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	var buildDuration *durationpb.Duration
	if event.BuildDuration != 0 {
		buildDuration = durationpb.New(event.BuildDuration)
	}

	msg := &eventsv1.ShipAssembled{
		EventUuid:     event.UUID,
		OrderUuid:     event.OrderUUID,
		UserUuid:      event.UserUUID,
		BuildDuration: buildDuration,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	err = p.kafkaProducer.Send(ctx, []byte(event.UUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
