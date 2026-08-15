package redis

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
	Password       string
	Logger         Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "redis-container",
		ImageName:     "redis:7-alpine",
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
