package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderServiceEnvConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
}

type orderServiceConfig struct {
	raw orderServiceEnvConfig
}

func NewOrderServiceConfig() (*orderServiceConfig, error) {
	var raw orderServiceEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderServiceConfig{raw: raw}, nil
}

func (cfg *orderServiceConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
