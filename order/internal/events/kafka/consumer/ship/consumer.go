package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	def "github.com/Reensef/go-microservices-course/order/internal/events"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	service "github.com/Reensef/go-microservices-course/order/internal/service"
	"github.com/Reensef/go-microservices-course/platform/pkg/kafka"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

var _ def.ShipConsumer = (*consumer)(nil)

type consumer struct {
	orderService service.OrderService
	group        sarama.ConsumerGroup
	middlewares  []kafka.Middleware
	topics       []string
}

func NewConsumer(
	orderService service.OrderService,
) *consumer {
	return &consumer{
		orderService: orderService,
	}
}

func (c *consumer) RunConsumer(ctx context.Context) error {
	err := c.Consume(ctx, c.shipHandler)
	if err != nil {
		logger.Error(ctx, "consume from order.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}

func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	newGroupHandler := kafka.NewGroupHandler(handler, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}

// Close останавливает консьюмер, закрывая группу консьюмеров.
// После вызова Close() выполняющийся Consume завершится с ошибкой sarama.ErrClosedConsumerGroup,
// которая обрабатывается как штатное завершение работы.
func (c *consumer) Close() error {
	return c.group.Close()
}

func (c *consumer) shipHandler(ctx context.Context, msg kafka.Message) error {
	event, err := c.shipDecode(msg.Value)
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
		zap.Duration("build_duration", event.BuildDuration),
	)

	err = c.orderService.AssembleOrder(ctx, event.OrderUUID)
	if err != nil {
		return fmt.Errorf("failed to mark order as assemble: %w", err)
	}

	return nil
}

func (c *consumer) shipDecode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		UUID:          pb.EventUuid,
		OrderUUID:     pb.OrderUuid,
		UserUUID:      pb.UserUuid,
		BuildDuration: pb.BuildDuration.AsDuration(),
	}, nil
}
