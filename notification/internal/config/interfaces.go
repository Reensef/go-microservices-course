package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
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
