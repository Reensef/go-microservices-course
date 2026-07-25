package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipProducerEnvConfig struct {
	Topic string `env:"SHIP_TOPIC_NAME,required"`
}

type shipProducerConfig struct {
	raw shipProducerEnvConfig
}

func NewShipProducerConfig() (*shipProducerConfig, error) {
	var raw shipProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipProducerConfig{raw: raw}, nil
}

func (cfg *shipProducerConfig) Topic() string {
	return cfg.raw.Topic
}

// Config возвращает конфигурацию для sarama sync producer
func (cfg *shipProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
