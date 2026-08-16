package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap/zapcore"

	"github.com/Reensef/go-microservices-course/notification/internal/config"
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
	logger.Info("notification service started")

	eg, egCtx := errgroup.WithContext(ctx)

	notificationConsumer, err := a.diContainer.NotificationConsumer(egCtx)
	if err != nil {
		return err
	}

	telegramBot, err := a.diContainer.TelegramBot(egCtx)
	if err != nil {
		return err
	}

	eg.Go(func() error {
		return notificationConsumer.RunConsumer(egCtx)
	})

	eg.Go(func() error {
		telegramBot.Start(egCtx)
		return nil
	})

	return eg.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
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
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(config.AppConfig().Logger.Level()))
	if err != nil {
		return err
	}

	opts := []logger.Option{
		logger.WithJSON(config.AppConfig().Logger.AsJson()),
	}

	if config.AppConfig().Logger.EnableOTLP() {
		opts = append(opts, logger.WithOTLP(config.AppConfig().Logger.OTLPEndpoint(), "notification-service", "dev"))
	}

	err = logger.Init(level, opts...)
	if err != nil {
		return err
	}

	closer.AddNamed("Logger OTLP", func(ctx context.Context) error {
		return logger.Close(ctx)
	})

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
