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

	appConfig = &config{
		Logger:           loggerConfig,
		InventoryService: inventoryServiceConfig,
		Mongo:            mongoConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
