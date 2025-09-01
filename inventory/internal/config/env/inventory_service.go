package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryServiceEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryServiceConfig struct {
	raw inventoryServiceEnvConfig
}

func NewInventoryServiceConfig() (*inventoryServiceConfig, error) {
	var raw inventoryServiceEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryServiceConfig{raw: raw}, nil
}

func (cfg *inventoryServiceConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
