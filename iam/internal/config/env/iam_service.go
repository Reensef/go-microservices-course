package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamServiceEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type iamServiceConfig struct {
	raw iamServiceEnvConfig
}

func NewIamServiceConfig() (*iamServiceConfig, error) {
	var raw iamServiceEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamServiceConfig{raw: raw}, nil
}

func (cfg *iamServiceConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
