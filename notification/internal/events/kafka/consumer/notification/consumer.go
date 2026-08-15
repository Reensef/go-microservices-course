package notification

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	def "github.com/Reensef/go-microservices-course/notification/internal/events"
	"github.com/Reensef/go-microservices-course/notification/internal/model"
	service "github.com/Reensef/go-microservices-course/notification/internal/service"
	"github.com/Reensef/go-microservices-course/platform/pkg/kafka"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
)

var _ def.Consumer = (*consumer)(nil)

type consumer struct {
	notificationService service.NotificationService

	group sarama.ConsumerGroup

	orderPaidTopic     string
	shipAssembledTopic string
}

func NewConsumer(
	brokers []string,
	groupID string,
	orderPaidTopic string,
	shipAssembledTopic string,
	saramaConfig *sarama.Config,
	notificationService service.NotificationService,
) (*consumer, error) {
	group, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &consumer{
		notificationService: notificationService,
		group:               group,
		orderPaidTopic:      orderPaidTopic,
		shipAssembledTopic:  shipAssembledTopic,
	}, nil
}

func (c *consumer) RunConsumer(ctx context.Context) error {
	logger.Info("Starting notification consumer service")

	topics := []string{c.orderPaidTopic, c.shipAssembledTopic}
	groupHandler := kafka.NewGroupHandler(c.handleMessage)

	for {
		if err := c.group.Consume(ctx, topics, groupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			logger.Error("consume from notification topics error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.Info("Kafka consumer group rebalancing...")
	}
}

// Close останавливает консьюмер, закрывая группу консьюмеров.
func (c *consumer) Close() error {
	return c.group.Close()
}

func (c *consumer) handleMessage(ctx context.Context, msg kafka.Message) error {
	switch msg.Topic {
	case c.orderPaidTopic:
		return c.handleOrderPaid(ctx, msg)
	case c.shipAssembledTopic:
		return c.handleShipAssembled(ctx, msg)
	default:
		logger.Warn("received message from unknown topic", zap.String("topic", msg.Topic))
		return nil
	}
}

func (c *consumer) handleOrderPaid(ctx context.Context, msg kafka.Message) error {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(msg.Value, &pb); err != nil {
		logger.Error("failed to unmarshal OrderPaid", zap.Error(err))
		return fmt.Errorf("failed to unmarshal OrderPaid: %w", err)
	}

	event := model.OrderPaidEvent{
		UUID:            pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		TransactionUUID: pb.TransactionUuid,
		PaymentMethod:   pb.PaymentMethod,
	}

	logger.Info("Processing OrderPaid message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
	)

	if err := c.notificationService.NotifyOrderPaid(ctx, event); err != nil {
		logger.Error("failed to notify OrderPaid", zap.Error(err))
		return fmt.Errorf("failed to notify OrderPaid: %w", err)
	}

	return nil
}

func (c *consumer) handleShipAssembled(ctx context.Context, msg kafka.Message) error {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(msg.Value, &pb); err != nil {
		logger.Error("failed to unmarshal ShipAssembled", zap.Error(err))
		return fmt.Errorf("failed to unmarshal ShipAssembled: %w", err)
	}

	event := model.ShipAssembledEvent{
		UUID:      pb.EventUuid,
		OrderUUID: pb.OrderUuid,
		UserUUID:  pb.UserUuid,
	}
	if pb.BuildDuration != nil {
		event.BuildDuration = pb.BuildDuration.AsDuration()
	}

	logger.Info("Processing ShipAssembled message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
	)

	if err := c.notificationService.NotifyShipAssembled(ctx, event); err != nil {
		logger.Error("failed to notify ShipAssembled", zap.Error(err))
		return fmt.Errorf("failed to notify ShipAssembled: %w", err)
	}

	return nil
}
