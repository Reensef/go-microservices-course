package redis

import (
	"go.uber.org/zap"

	tc "github.com/Reensef/go-microservices-course/platform/pkg/testcontainers"
)

type Logger = tc.Logger

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
