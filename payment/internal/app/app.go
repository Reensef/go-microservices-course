package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Reensef/go-microservices-course/payment/internal/config"
	"github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/grpc/health"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
	paymentProtoApi "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initTracing,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(config.AppConfig().Logger.Level()))
	if err != nil {
		return err
	}

	opts := []logger.Option{logger.WithJSON(config.AppConfig().Logger.AsJson())}
	if config.AppConfig().Logger.EnableOTLP() {
		opts = append(opts, logger.WithOTLP(
			config.AppConfig().Logger.OTLPEndpoint(),
			config.AppConfig().Service.Name(),
			config.AppConfig().Service.Environment(),
		))
	}
	err = logger.Init(ctx, level, opts...)
	if err != nil {
		return err
	}

	closer.AddNamed("Logger OTLP", logger.Close)

	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	cfg := config.AppConfig().Tracing
	service := config.AppConfig().Service

	err := tracer.Init(ctx,
		cfg.CollectorEndpoint(),
		service.Name(),
		service.Environment(),
		tracer.WithServiceVersion(cfg.ServiceVersion()),
		tracer.WithInsecure(),
	)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracer.Shutdown)

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().PaymentService.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			tracer.UnaryServerInterceptor(),
		),
	)
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	paymentProtoApi.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentApi(ctx))

	return nil
}

func (a *App) runGRPCServer(_ context.Context) error {
	logger.Info(fmt.Sprintf(
		"🚀 gRPC PaymentService server listening on %s",
		config.AppConfig().PaymentService.Address(),
	))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
