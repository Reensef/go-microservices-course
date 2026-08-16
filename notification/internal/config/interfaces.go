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

type NotificationConsumerConfig interface {
	GroupID() string
	OrderPaidTopic() string
	ShipAssembledTopic() string
	Config() *sarama.Config
}

type TelegramConfig interface {
	BotToken() string
	ChatID() int64
}

type ServiceConfig interface {
	Name() string
	Environment() string
}
