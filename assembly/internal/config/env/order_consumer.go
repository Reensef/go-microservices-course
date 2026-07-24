package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderConsumerEnvConfig struct {
	GroupID string `env:"ASSEMBLY_CONSUMER_GROUP_ID,required"`
	Topic   string `env:"ORDER_TOPIC_NAME,required"`
}

type orderConsumerConfig struct {
	raw orderConsumerEnvConfig
}

func NewOrderConsumerConfig() (*orderConsumerConfig, error) {
	var raw orderConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderConsumerConfig{raw: raw}, nil
}

func (cfg *orderConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

// Config возвращает конфигурацию для sarama consumer group
func (cfg *orderConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
