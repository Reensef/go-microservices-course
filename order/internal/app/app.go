package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/Reensef/go-microservices-course/order/internal/config"
	"github.com/Reensef/go-microservices-course/order/internal/metric"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	"github.com/Reensef/go-microservices-course/platform/pkg/tracer"
)

type App struct {
	diContainer     *diContainer
	orderHttpServer *http.Server
	orderRouter     *chi.Mux
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
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return a.runOrderHttpServer(egCtx)
	})

	eg.Go(func() error {
		return a.diContainer.ShipConsumer(egCtx).RunConsumer(egCtx)
	})

	return eg.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initTracing,
		a.applyMigrations,
		a.initCloser,
		a.initMetrics,
		a.initOrderRouter,
		a.initOrderHttpServer,
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

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) applyMigrations(ctx context.Context) error {
	return a.diContainer.SqlMigrator(ctx).Up()
}

func (a *App) initOrderRouter(ctx context.Context) error {
	a.orderRouter = chi.NewRouter()

	a.orderRouter.Use(middleware.Logger)
	a.orderRouter.Use(middleware.Recoverer)
	a.orderRouter.Use(middleware.Timeout(10 * time.Second))

	a.orderRouter.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	a.orderRouter.Group(func(r chi.Router) {
		r.Use(a.diContainer.AuthMiddleware(ctx))
		r.Mount("/", a.diContainer.OrderApi(ctx))
	})

	return nil
}

func (a *App) initOrderHttpServer(ctx context.Context) error {
	a.orderHttpServer = &http.Server{
		Addr:              config.AppConfig().OrderService.Address(),
		Handler:           a.orderRouter,
		ReadHeaderTimeout: 5 * time.Second, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	return nil
}

func (a *App) runOrderHttpServer(_ context.Context) error {
	logger.Info(fmt.Sprintf(
		"🚀 HTTP Order Service server listening on %s",
		config.AppConfig().OrderService.Address()),
	)
	closer.AddNamed("Order HTTP server", func(ctx context.Context) error {
		err := a.orderHttpServer.Shutdown(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Warn("Order HTTP server shutdown error", zap.Error(err))
			return err
		}

		return nil
	})

	err := a.orderHttpServer.ListenAndServe()
	if err != nil {
		logger.Warn(fmt.Sprintf("Order HTTP server startup error: %v", err))
		return err
	}

	return nil
}
