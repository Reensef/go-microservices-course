package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type notificationConsumerEnvConfig struct {
	GroupID            string `env:"CONSUMER_GROUP_ID,required"`
	OrderPaidTopic     string `env:"ORDER_PAID_TOPIC_NAME,required"`
	ShipAssembledTopic string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
}

type notificationConsumerConfig struct {
	raw notificationConsumerEnvConfig
}

func NewNotificationConsumerConfig() (*notificationConsumerConfig, error) {
	var raw notificationConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &notificationConsumerConfig{raw: raw}, nil
}

func (cfg *notificationConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *notificationConsumerConfig) OrderPaidTopic() string {
	return cfg.raw.OrderPaidTopic
}

func (cfg *notificationConsumerConfig) ShipAssembledTopic() string {
	return cfg.raw.ShipAssembledTopic
}

// Config возвращает конфигурацию для sarama consumer group
func (cfg *notificationConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
