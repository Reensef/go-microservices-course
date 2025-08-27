package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderHandler "github.com/Reensef/go-microservices-course/order/internal/api/order/v1"
	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	inventoryClient "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Reensef/go-microservices-course/order/internal/client/grpc/payment/v1"
	"github.com/Reensef/go-microservices-course/order/internal/config"
	repo "github.com/Reensef/go-microservices-course/order/internal/repository"
	orderRepo "github.com/Reensef/go-microservices-course/order/internal/repository/order"
	service "github.com/Reensef/go-microservices-course/order/internal/service"
	orderService "github.com/Reensef/go-microservices-course/order/internal/service/order"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/sqlmigrator"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	paymentGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderHandler orderApi.Handler
	orderApi     *orderApi.Server
	orderService service.OrderService
	orderRepo    repo.OrderRepository

	inventoryClient grpcClients.IntentoryClient
	paymentClient   grpcClients.PaymentClient

	inventoryGrpc inventoryGrpc.InventoryServiceClient
	paymentGrpc   paymentGrpc.PaymentServiceClient

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
