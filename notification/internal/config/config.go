package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Reensef/go-microservices-course/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig

	Kafka                KafkaConfig
	NotificationConsumer NotificationConsumerConfig

	Telegram TelegramConfig
	Service  ServiceConfig
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

	notificationConsumerCfg, err := env.NewNotificationConsumerConfig()
	if err != nil {
		return err
	}

	telegramCfg, err := env.NewTelegramConfig()
	if err != nil {
		return err
	}

	serviceCfg, err := env.NewServiceConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:               loggerCfg,
		Kafka:                kafkaCfg,
		NotificationConsumer: notificationConsumerCfg,
		Telegram:             telegramCfg,
		Service:              serviceCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
