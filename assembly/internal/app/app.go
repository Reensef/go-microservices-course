package app

import (
	"context"

	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/Reensef/go-microservices-course/assembly/internal/config"
	"github.com/Reensef/go-microservices-course/assembly/internal/metric"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
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
	logger.Info("🚀 assembly service started")

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return a.diContainer.OrderConsumer(egCtx).RunConsumer(egCtx)
	})

	return eg.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMetrics,
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

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	meterProvider, err := metric.Init(ctx, config.AppConfig().Metrics.OTLPEndpoint(), config.AppConfig().Service.Name())
	if err != nil {
		return err
	}

	closer.AddNamed("OTel MeterProvider", func(ctx context.Context) error {
		return meterProvider.Shutdown(ctx)
	})

	return nil
}
