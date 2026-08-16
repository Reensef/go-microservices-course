package env

import "github.com/caarlos0/env/v11"

type serviceEnvConfig struct {
	NameValue        string `env:"OTEL_SERVICE_NAME,required"`
	EnvironmentValue string `env:"OTEL_ENVIRONMENT,required"`
}

type serviceConfig struct {
	raw serviceEnvConfig
}

func NewServiceConfig() (*serviceConfig, error) {
	var raw serviceEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &serviceConfig{raw: raw}, nil
}

func (cfg *serviceConfig) Name() string {
	return cfg.raw.NameValue
}

func (cfg *serviceConfig) Environment() string {
	return cfg.raw.EnvironmentValue
}
