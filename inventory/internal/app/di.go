package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1"
	"github.com/Reensef/go-microservices-course/inventory/internal/config"
	repository "github.com/Reensef/go-microservices-course/inventory/internal/repository"
	partRepository "github.com/Reensef/go-microservices-course/inventory/internal/repository/part"
	service "github.com/Reensef/go-microservices-course/inventory/internal/service"
	inventoryService "github.com/Reensef/go-microservices-course/inventory/internal/service/inventory"
	"github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/grpc/health"
	inventoryProtoApi "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryApi        inventoryProtoApi.InventoryServiceServer
	inventoryGrpcServer *grpc.Server
	inventoryListener   net.Listener

	inventoryService service.InventoryService

	partRepository repository.PartRepository
	mongoHandler   *mongo.Database
	mongoClient    *mongo.Client
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryApi(ctx context.Context) inventoryProtoApi.InventoryServiceServer {
	if d.inventoryApi == nil {
		d.inventoryApi = inventoryApi.New(d.InventoryService(ctx))
	}

	return d.inventoryApi
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.New(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryGrpcServer(ctx context.Context) *grpc.Server {
	if d.inventoryGrpcServer == nil {
		grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

		reflection.Register(grpcServer)

		// Регистрируем health service для проверки работоспособности
		health.RegisterService(grpcServer)

		inventoryProtoApi.RegisterInventoryServiceServer(grpcServer, d.InventoryApi(ctx))

		d.inventoryGrpcServer = grpcServer
	}

	return d.inventoryGrpcServer
}

func (d *diContainer) InventoryListener(ctx context.Context) net.Listener {
	if d.inventoryListener == nil {
		listener, err := net.Listen("tcp", config.AppConfig().InventoryService.Address())
		if err != nil {
			panic(fmt.Errorf("failed to listen: %w", err))
		}

		closer.AddNamed("TCP listener", func(ctx context.Context) error {
			lerr := listener.Close()
			if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
				return lerr
			}

			return nil
		})

		d.inventoryListener = listener
	}

	return d.inventoryListener
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.New(d.MongoHandler(ctx))
	}

	return d.partRepository
}

func (d *diContainer) MongoHandler(ctx context.Context) *mongo.Database {
	if d.mongoHandler == nil {
		d.mongoHandler = d.MongoClient(ctx).Database(
			config.AppConfig().Mongo.DatabaseName(),
		)
	}

	return d.mongoHandler
}

func (d *diContainer) MongoClient(ctx context.Context) *mongo.Client {
	if d.mongoClient == nil {
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(
			config.AppConfig().Mongo.URI(),
		))
		if err != nil {
			panic(fmt.Errorf("failed to connect to mongo: %w", err))
		}

		closer.AddNamed("Mongo client", func(ctx context.Context) error {
			return mongoClient.Disconnect(ctx)
		})

		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			panic(fmt.Errorf("failed to ping database: %w", err))
		}

		d.mongoClient = mongoClient
	}

	return d.mongoClient
}
