package ship

import (
	"context"

	"github.com/IBM/sarama"
	def "github.com/Reensef/go-microservices-course/order/internal/events"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var _ def.OrderProducer = (*producer)(nil)

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

func (p *producer) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsv1.OrderPaid{
		EventUuid:       event.UUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		TransactionUuid: event.TransactionUUID,
		PaymentMethod:   event.PaymentMethod.String(),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder([]byte(event.UUID)),
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return err
	}

	return nil
}
