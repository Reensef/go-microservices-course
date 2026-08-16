package mongo

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
	AuthDB         string
	Logger         Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "mongo-container",
		ImageName:     "mongo:8.0",
		Database:      "test",
		Username:      "root",
		Password:      "root",
		AuthDB:        "admin",
		Logger:        zap.NewNop(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
