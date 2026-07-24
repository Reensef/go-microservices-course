package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type kafkaEnvConfig struct {
	Host string `env:"KAFKA_HOST,required"`
	Port string `env:"KAFKA_PORT,required"`
}

type kafkaConfig struct {
	raw kafkaEnvConfig
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var raw kafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{raw: raw}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return []string{net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)}
}
