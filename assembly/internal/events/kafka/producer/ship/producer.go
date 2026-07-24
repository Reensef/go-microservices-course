package ship

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	events "github.com/Reensef/go-microservices-course/assembly/internal/events"
	"github.com/Reensef/go-microservices-course/assembly/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
)

var _ events.ShipProducer = (*producer)(nil)

// SaramaSyncProducer is an interface for mock sarama sync producer.
type SaramaSyncProducer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type producer struct {
	syncProducer SaramaSyncProducer
	topic        string
}

func NewProducer(syncProducer SaramaSyncProducer, topic string) *producer {
	return &producer{
		syncProducer: syncProducer,
		topic:        topic,
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

	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder([]byte(event.UUID)),
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
