package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentServiceEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentServiceConfig struct {
	raw paymentServiceEnvConfig
}

func NewPaymentServiceConfig() (*paymentServiceConfig, error) {
	var raw paymentServiceEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentServiceConfig{raw: raw}, nil
}

func (cfg *paymentServiceConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
