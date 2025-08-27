package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Reensef/go-microservices-course/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig

	OrderService    OrderServiceConfig
	PaymentClient   PaymentClientConfig
	InventoryClient InventoryClientConfig

	SqlMigrator SqlMigratorConfig

	Postgres PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderService, err := env.NewOrderServiceConfig()
	if err != nil {
		return err
	}

	inventoryClient, err := env.NewInventoryClientConfig()
	if err != nil {
		return err
	}

	paymentClient, err := env.NewPaymentClientConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	sqlMigratorCfg, err := env.NewSqlMigratorConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:          loggerCfg,
		OrderService:    orderService,
		PaymentClient:   paymentClient,
		InventoryClient: inventoryClient,
		Postgres:        postgresCfg,
		SqlMigrator:     sqlMigratorCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
