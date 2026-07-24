package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipConsumerEnvConfig struct {
	GroupID string `env:"ORDER_CONSUMER_GROUP_ID,required"`
	Topic   string `env:"SHIP_TOPIC_NAME,required"`
}

type shipConsumerConfig struct {
	raw shipConsumerEnvConfig
}

func NewShipConsumerConfig() (*shipConsumerConfig, error) {
	var raw shipConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipConsumerConfig{raw: raw}, nil
}

func (cfg *shipConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *shipConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

// Config возвращает конфигурацию для sarama consumer group
func (cfg *shipConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
