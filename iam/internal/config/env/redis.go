package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     string `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type redisConfig struct {
	raw redisEnvConfig
}

func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &redisConfig{raw: raw}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *redisConfig) Password() string {
	return cfg.raw.Password
}

func (cfg *redisConfig) DB() int {
	return cfg.raw.DB
}
