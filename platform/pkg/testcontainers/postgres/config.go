package postgres

import (
	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
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
		Logger:        zap.NewNop(),
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
