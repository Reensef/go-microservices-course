package app

import (
	"context"
	"fmt"

	"github.com/Reensef/go-microservices-course/inventory/internal/config"
	"github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
	"go.uber.org/zap/zapcore"
)

type App struct {
	di *diContainer
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
	a.di = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	var level zapcore.Level
	_ = level.UnmarshalText([]byte(config.AppConfig().Logger.Level()))
	opts := []logger.Option{logger.WithJSON(config.AppConfig().Logger.AsJson())}
	if endpoint := config.AppConfig().Logger.OTLPEndpoint(); endpoint != "" {
		opts = append(opts, logger.WithOTLP(endpoint, "inventory-service", "dev"))
	}
	return logger.Init(level, opts...)
}

func (a *App) initTracing(ctx context.Context) error {
	cfg := config.AppConfig().Tracing

	if err := tracer.Init(ctx,
		cfg.CollectorEndpoint(),
		cfg.ServiceName(),
		cfg.Environment(),
		tracer.WithServiceVersion(cfg.ServiceVersion()),
		tracer.WithInsecure(),
	); err != nil {
		return err
	}

	closer.AddNamed("tracer", tracer.Shutdown)

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	closer.AddNamed("Logger OTLP", func(ctx context.Context) error {
		return logger.Close(ctx)
	})
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(fmt.Sprintf(
		"🚀 gRPC InventoryService server listening on %s",
		a.di.InventoryListener(ctx).Addr(),
	))

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.di.InventoryGrpcServer(ctx).GracefulStop()
		return nil
	})

	err := a.di.InventoryGrpcServer(ctx).Serve(a.di.InventoryListener(ctx))
	if err != nil {
		return err
	}

	return nil
}
