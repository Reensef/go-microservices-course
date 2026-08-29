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

	Kafka         KafkaConfig
	OrderProducer OrderProducerConfig
	ShipConsumer  ShipConsumerConfig

	SqlMigrator SqlMigratorConfig

	Postgres PostgresConfig
	Metrics  MetricsConfig
	Tracing  TracingConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderProducerCfg, err := env.NewOrderProducerConfig()
	if err != nil {
		return err
	}

	shipConsumerCfg, err := env.NewShipConsumerConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := env.NewMetricsConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	serviceCfg, err := env.NewServiceConfig()
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
		Kafka:           kafkaCfg,
		OrderProducer:   orderProducerCfg,
		ShipConsumer:    shipConsumerCfg,
		Metrics:         metricsCfg,
		Tracing:         tracingCfg,
		Service:         serviceCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
