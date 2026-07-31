package postgres

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName    string
	NetworkAliases []string
	ContainerName  string
	ImageName      string
	Database       string
	Username       string
	Password       string
	Logger         Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "postgres-container",
		ImageName:     "postgres:16-alpine",
		Database:      "test",
		Username:      "postgres",
		Password:      "postgres",
		Logger:        &logger.DummyLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
