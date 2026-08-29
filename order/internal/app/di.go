package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderMiddleware "github.com/Reensef/go-microservices-course/order/internal/api/middleware"
	orderHandler "github.com/Reensef/go-microservices-course/order/internal/api/order/v1"
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	inventoryClient "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Reensef/go-microservices-course/order/internal/client/grpc/payment/v1"
	"github.com/Reensef/go-microservices-course/order/internal/config"
	events "github.com/Reensef/go-microservices-course/order/internal/events"
	shipconsumer "github.com/Reensef/go-microservices-course/order/internal/events/kafka/consumer/ship"
	orderproducer "github.com/Reensef/go-microservices-course/order/internal/events/kafka/producer/order"
	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
	orderRepo "github.com/Reensef/go-microservices-course/order/internal/repository/order"
	service "github.com/Reensef/go-microservices-course/order/internal/service"
	orderService "github.com/Reensef/go-microservices-course/order/internal/service/order"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/sqlmigrator"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	paymentGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderHandler       orderApi.Handler
	orderApi           *orderApi.Server
	orderService       service.OrderService
	orderRepo          repo.OrderRepository
	orderProducer      events.OrderProducer
	saramaSyncProducer sarama.SyncProducer
	shipConsumer       events.ShipConsumer
	consumerGroup      sarama.ConsumerGroup
	inventoryClient    grpcClients.IntentoryClient
	paymentClient      grpcClients.PaymentClient

	inventoryGrpc inventoryGrpc.InventoryServiceClient
	paymentGrpc   paymentGrpc.PaymentServiceClient

	authMiddleware func(http.Handler) http.Handler

	postgresPool *pgxpool.Pool
	sqlMigrator  *sqlmigrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderApi(ctx context.Context) *orderApi.Server {
	if d.orderApi == nil {
		var err error
		d.orderApi, err = orderApi.NewServer(d.OrderHandler(ctx))
		if err != nil {
			panic(err)
		}
	}

	return d.orderApi
}

func (d *diContainer) OrderHandler(ctx context.Context) orderApi.Handler {
	if d.orderHandler == nil {
		d.orderHandler = orderHandler.NewHandler(d.OrderService(ctx))
	}

	return d.orderHandler
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.New(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderProducer(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repo.OrderRepository {
	if d.orderRepo == nil {
		d.orderRepo = orderRepo.New(d.PostgresPool(ctx))
	}

	return d.orderRepo
}

func (d *diContainer) OrderProducer(ctx context.Context) events.OrderProducer {
	if d.orderProducer == nil {
		d.orderProducer = orderproducer.NewProducer(
			d.OrderSyncProducer(ctx),
			config.AppConfig().OrderProducer.Topic(),
		)
	}

	return d.orderProducer
}

func (d *diContainer) OrderSyncProducer(ctx context.Context) sarama.SyncProducer {
	if d.saramaSyncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.saramaSyncProducer = p
	}

	return d.saramaSyncProducer
}

func (d *diContainer) ShipConsumer(ctx context.Context) events.ShipConsumer {
	if d.shipConsumer == nil {
		d.shipConsumer = shipconsumer.NewConsumer(
			d.OrderService(ctx),
			d.ConsumerGroup(ctx),
			config.AppConfig().ShipConsumer.Topic(),
		)
	}

	return d.shipConsumer
}

func (d *diContainer) ConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipConsumer.GroupID(),
			config.AppConfig().ShipConsumer.Config(),
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

func (d *diContainer) SqlMigrator(ctx context.Context) *sqlmigrator.Migrator {
	if d.sqlMigrator == nil {
		postgresPool := d.PostgresPool(ctx)
		d.sqlMigrator = sqlmigrator.New(
			stdlib.OpenDB(*postgresPool.Config().ConnConfig.Copy()),
			config.AppConfig().SqlMigrator.MigrationsDir(),
		)
	}

	return d.sqlMigrator
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		var err error
		d.postgresPool, err = pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres: %w", err))
		}

		err = d.postgresPool.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping postgres: %w", err))
		}

		closer.AddNamed("Postgres pool", func(ctx context.Context) error {
			d.postgresPool.Close()
			return nil
		})
	}

	return d.postgresPool
}

func (d *diContainer) InventoryClient(ctx context.Context) grpcClients.IntentoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryClient.New(d.InventoryGrpc(ctx))
	}

	return d.inventoryClient
}

func (d *diContainer) InventoryGrpc(ctx context.Context) inventoryGrpc.InventoryServiceClient {
	if d.inventoryGrpc == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				tracer.UnaryClientInterceptor(),
			),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory service: %v\n", err))
		}

		closer.AddNamed("Inventory gRPC client", func(ctx context.Context) error {
			if err := conn.Close(); err != nil {
				return err
			}
			return nil
		})

		d.inventoryGrpc = inventoryGrpc.NewInventoryServiceClient(conn)
	}

	return d.inventoryGrpc
}

func (d *diContainer) PaymentClient(ctx context.Context) grpcClients.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentClient.New(d.PaymentGrpc(ctx))
	}

	return d.paymentClient
}

func (d *diContainer) PaymentGrpc(ctx context.Context) paymentGrpc.PaymentServiceClient {
	if d.paymentGrpc == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(
				tracer.UnaryClientInterceptor(),
			),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment service: %v\n", err))
		}

		closer.AddNamed("Payment gRPC client", func(ctx context.Context) error {
			if err := conn.Close(); err != nil {
				return err
			}
			return nil
		})

		d.paymentGrpc = paymentGrpc.NewPaymentServiceClient(conn)
	}

	return d.paymentGrpc
}

func (d *diContainer) AuthMiddleware(ctx context.Context) func(http.Handler) http.Handler {
	if d.authMiddleware == nil {
		d.authMiddleware = orderMiddleware.NewAuthMiddleware()
	}

	return d.authMiddleware
}
