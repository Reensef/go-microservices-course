package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
}

type OrderServiceConfig interface {
	Address() string
}

type PaymentClientConfig interface {
	Address() string
}

type InventoryClientConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type PostgresConfig interface {
	URI() string
}

type SqlMigratorConfig interface {
	MigrationsDir() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type ShipConsumerConfig interface {
	GroupID() string
	Topic() string
	Config() *sarama.Config
}

type MetricsConfig interface {
	OTLPEndpoint() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceVersion() string
}

type ServiceConfig interface {
	Name() string
	Environment() string
}
