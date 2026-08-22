package app

import (
	"context"
	"fmt"
	"net"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	authApi "github.com/Reensef/go-microservices-course/iam/internal/api/auth/v1"
	authzApi "github.com/Reensef/go-microservices-course/iam/internal/api/authz/v1"
	userApi "github.com/Reensef/go-microservices-course/iam/internal/api/user/v1"
	"github.com/Reensef/go-microservices-course/iam/internal/config"
	repository "github.com/Reensef/go-microservices-course/iam/internal/repository"
	sessionRepository "github.com/Reensef/go-microservices-course/iam/internal/repository/session"
	userRepository "github.com/Reensef/go-microservices-course/iam/internal/repository/user"
	service "github.com/Reensef/go-microservices-course/iam/internal/service"
	authService "github.com/Reensef/go-microservices-course/iam/internal/service/auth"
	userService "github.com/Reensef/go-microservices-course/iam/internal/service/user"
	"github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/grpc/health"
	"github.com/Reensef/go-microservices-course/platform/pkg/sqlmigrator"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

type diContainer struct {
	authApi  iamV1.AuthServiceServer
	userApi  iamV1.UserServiceServer
	authzApi authv3.AuthorizationServer

	grpcServer *grpc.Server
	listener   net.Listener

	authService service.AuthService
	userService service.UserService

	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository

	postgresPool *pgxpool.Pool
	sqlMigrator  *sqlmigrator.Migrator
	redisClient  *redis.Client
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthApi(ctx context.Context) iamV1.AuthServiceServer {
	if d.authApi == nil {
		d.authApi = authApi.New(d.AuthService(ctx))
	}

	return d.authApi
}

func (d *diContainer) UserApi(ctx context.Context) iamV1.UserServiceServer {
	if d.userApi == nil {
		d.userApi = userApi.New(d.UserService(ctx))
	}

	return d.userApi
}

func (d *diContainer) AuthzApi(ctx context.Context) authv3.AuthorizationServer {
	if d.authzApi == nil {
		d.authzApi = authzApi.New(d.AuthService(ctx))
	}

	return d.authzApi
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.New(
			d.UserRepository(ctx),
			d.SessionRepository(ctx),
			config.AppConfig().Session.TTL(),
		)
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.New(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepo == nil {
		d.userRepo = userRepository.New(d.PostgresPool(ctx))
	}

	return d.userRepo
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepo == nil {
		d.sessionRepo = sessionRepository.New(d.RedisClient(ctx))
	}

	return d.sessionRepo
}

func (d *diContainer) GrpcServer(ctx context.Context) *grpc.Server {
	if d.grpcServer == nil {
		grpcServer := grpc.NewServer(
			grpc.Creds(insecure.NewCredentials()),
			grpc.ChainUnaryInterceptor(tracer.UnaryServerInterceptor()),
		)

		reflection.Register(grpcServer)

		// Регистрируем health service для проверки работоспособности
		health.RegisterService(grpcServer)

		iamV1.RegisterAuthServiceServer(grpcServer, d.AuthApi(ctx))
		iamV1.RegisterUserServiceServer(grpcServer, d.UserApi(ctx))
		authv3.RegisterAuthorizationServer(grpcServer, d.AuthzApi(ctx))

		d.grpcServer = grpcServer
	}

	return d.grpcServer
}

func (d *diContainer) Listener(_ context.Context) net.Listener {
	if d.listener == nil {
		listener, err := net.Listen("tcp", config.AppConfig().IamService.Address())
		if err != nil {
			panic(fmt.Errorf("failed to listen: %w", err))
		}

		closer.AddNamed("TCP listener", func(ctx context.Context) error {
			return listener.Close()
		})

		d.listener = listener
	}

	return d.listener
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

func (d *diContainer) RedisClient(ctx context.Context) *redis.Client {
	if d.redisClient == nil {
		client := redis.NewClient(&redis.Options{
			Addr:     config.AppConfig().Redis.Address(),
			Password: config.AppConfig().Redis.Password(),
			DB:       config.AppConfig().Redis.DB(),
		})

		if err := client.Ping(ctx).Err(); err != nil {
			panic(fmt.Errorf("failed to ping redis: %w", err))
		}

		closer.AddNamed("Redis client", func(ctx context.Context) error {
			return client.Close()
		})

		d.redisClient = client
	}

	return d.redisClient
}
