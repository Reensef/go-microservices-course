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

	"github.com/Reensef/go-microservices-course/order/internal/config"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
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
	return a.runOrderHttpServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.applyMigrations,
		a.initCloser,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
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

	a.orderRouter.Mount("/", a.diContainer.OrderApi(ctx))

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

func (a *App) runOrderHttpServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf(
		"🚀 HTTP Order Service server listening on %s",
		config.AppConfig().OrderService.Address()),
	)
	closer.AddNamed("Order HTTP server", func(ctx context.Context) error {
		err := a.orderHttpServer.Shutdown(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Warn(ctx, "Order HTTP server shutdown error", zap.Error(err))
			return err
		}

		return nil
	})

	err := a.orderHttpServer.ListenAndServe()
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("Order HTTP server startup error: %v", err))
		return err
	}

	return nil
}
