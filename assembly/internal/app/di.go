package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/Reensef/go-microservices-course/assembly/internal/config"
	events "github.com/Reensef/go-microservices-course/assembly/internal/events"
	orderConsumer "github.com/Reensef/go-microservices-course/assembly/internal/events/kafka/consumer/order"
	shipProducer "github.com/Reensef/go-microservices-course/assembly/internal/events/kafka/producer/ship"
	service "github.com/Reensef/go-microservices-course/assembly/internal/service"
	assemlyservice "github.com/Reensef/go-microservices-course/assembly/internal/service/assembly"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
)

type diContainer struct {
	// services
	assemblyService service.AssemblyService

	// adapters
	orderConsumer events.OrderConsumer
	shipProducer  events.ShipProducer

	// third-party clients
	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderConsumer(ctx context.Context) events.OrderConsumer {
	if d.orderConsumer == nil {
		d.orderConsumer = orderConsumer.NewConsumer(
			d.AssemblyService(ctx),
			d.ConsumerGroup(ctx),
			config.AppConfig().OrderConsumer.Topic(),
		)
	}

	return d.orderConsumer
}

func (d *diContainer) ShipProducer(ctx context.Context) events.ShipProducer {
	if d.shipProducer == nil {
		d.shipProducer = shipProducer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig().ShipProducer.Topic(),
		)
	}

	return d.shipProducer
}

func (d *diContainer) ConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderConsumer.GroupID(),
			config.AppConfig().OrderConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create kafka consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}

func (d *diContainer) SyncProducer(_ context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create kafka sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) AssemblyService(ctx context.Context) service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assemlyservice.NewService(
			d.ShipProducer(ctx),
		)
	}

	return d.assemblyService
}
