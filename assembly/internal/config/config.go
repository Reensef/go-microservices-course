package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Reensef/go-microservices-course/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig

	Kafka         KafkaConfig
	OrderConsumer OrderConsumerConfig
	ShipProducer  ShipProducerConfig
	Metrics       MetricsConfig
}

func Load(envFile string) error {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderConsumerCfg, err := env.NewOrderConsumerConfig()
	if err != nil {
		return err
	}

	shipProducerCfg, err := env.NewShipProducerConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := env.NewMetricsConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		Kafka:         kafkaCfg,
		OrderConsumer: orderConsumerCfg,
		ShipProducer:  shipProducerCfg,
		Metrics:       metricsCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
