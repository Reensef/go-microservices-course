package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderConsumerConfig interface {
	GroupID() string
	Topic() string
	Config() *sarama.Config
}

type ShipProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type MetricsConfig interface {
	OTLPEndpoint() string
}

type ServiceConfig interface {
	Name() string
	Environment() string
}
