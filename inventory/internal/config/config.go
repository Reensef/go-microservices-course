package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Reensef/go-microservices-course/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger           LoggerConfig
	InventoryService InventoryServiceConfig
	Mongo            MongoConfig
	Tracing          TracingConfig
	Service          ServiceConfig
}

func Load(envFile string) error {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryServiceConfig, err := env.NewInventoryServiceConfig()
	if err != nil {
		return err
	}

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	tracingConfig, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	serviceConfig, err := env.NewServiceConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:           loggerConfig,
		InventoryService: inventoryServiceConfig,
		Mongo:            mongoConfig,
		Tracing:          tracingConfig,
		Service:          serviceConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
