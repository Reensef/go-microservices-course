package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/Reensef/go-microservices-course/assembly/internal/config"
	"github.com/Reensef/go-microservices-course/assembly/internal/metric"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"go.uber.org/zap/zapcore"
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
	logger.Info(ctx, "🚀 assembly service started")

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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		[]zapcore.Core{
			logger.NewStdoutCore(logger.StdoutCoreConfig{
				Level:  config.AppConfig().Logger.Level(),
				AsJSON: config.AppConfig().Logger.AsJson(),
			}),
		},
		logger.DefaultZapOpts()...,
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	meterProvider, err := metric.Init(ctx, config.AppConfig().Metrics.OTLPEndpoint())
	if err != nil {
		return err
	}

	closer.AddNamed("OTel MeterProvider", func(ctx context.Context) error {
		return meterProvider.Shutdown(ctx)
	})

	return nil
}
