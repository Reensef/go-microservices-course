package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/inventory/internal/app"
	"github.com/Reensef/go-microservices-course/inventory/internal/config"
	"github.com/Reensef/go-microservices-course/platform/pkg/closer"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

func main() {
	var env string
	flag.StringVar(&env, "env", "", "Путь к файлу с переменными окружения")
	flag.Parse()

	err := config.Load(env)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error("error creating application", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error("error running application", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error("Error shutting down", zap.Error(err))
	}
}
