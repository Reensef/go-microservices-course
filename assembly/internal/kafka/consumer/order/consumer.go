package order

import (
	"context"
	"fmt"

	def "github.com/Reensef/go-microservices-course/assembly/internal/kafka"
	"github.com/Reensef/go-microservices-course/assembly/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/kafka"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

var _ def.OrderConsumer = (*consumer)(nil)

type consumer struct {
	kafkaConsumer kafka.Consumer
}

func NewConsumer(
	kafkaConsumer kafka.Consumer,
) *consumer {
	return &consumer{
		kafkaConsumer: kafkaConsumer,
	}
}

func (c *consumer) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order ufoRecordedConsumer service")

	err := c.kafkaConsumer.Consume(ctx, c.OrderHandler)
	if err != nil {
		logger.Error(ctx, "consume from order.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}

func (c *consumer) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := c.OrderDecode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode UFORecorded", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("transaction_uuid", event.TransactionUUID),
		zap.String("payment_method", event.PaymentMethod),
	)

	return nil
}

func (c *consumer) OrderDecode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		UUID:            pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		TransactionUUID: pb.TransactionUuid,
		PaymentMethod:   pb.PaymentMethod,
	}, nil
}
