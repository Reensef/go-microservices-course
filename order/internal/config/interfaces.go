package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
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
