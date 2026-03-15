package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Reensef/go-microservices-course/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger         LoggerConfig
	PaymentService PaymentServiceConfig
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

	paymentService, err := env.NewPaymentServiceConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:         loggerCfg,
		PaymentService: paymentService,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
