package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	events "github.com/Reensef/go-microservices-course/assembly/internal/events"
	"github.com/Reensef/go-microservices-course/assembly/internal/metric"
	"github.com/Reensef/go-microservices-course/assembly/internal/model"
	service "github.com/Reensef/go-microservices-course/assembly/internal/service"
	"github.com/Reensef/go-microservices-course/platform/pkg/kafka"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	eventsv1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/events/v1"
)

var _ events.OrderConsumer = (*consumer)(nil)

type consumer struct {
	group           sarama.ConsumerGroup
	topics          []string
	assemblyService service.AssemblyService
}

func NewConsumer(
	assemblyService service.AssemblyService,
	group sarama.ConsumerGroup,
	topic string,
) *consumer {
	return &consumer{
		group:           group,
		topics:          []string{topic},
		assemblyService: assemblyService,
	}
}

func (c *consumer) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")

	groupHandler := kafka.NewGroupHandler(c.OrderHandler)

	for {
		if err := c.group.Consume(ctx, c.topics, groupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			logger.Error(ctx, "consume from order topic error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}

// Close останавливает консьюмер, закрывая группу консьюмеров.
func (c *consumer) Close() error {
	return c.group.Close()
}

func (c *consumer) OrderHandler(ctx context.Context, msg kafka.Message) error {
	metric.IncOrdersReceived(ctx)

	event, err := c.OrderDecode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	err = c.assemblyService.AssembleShip(ctx, event.OrderUUID, event.UserUUID)
	if err != nil {
		return err
	}

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
