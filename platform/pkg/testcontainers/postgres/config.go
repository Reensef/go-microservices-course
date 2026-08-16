package postgres

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
