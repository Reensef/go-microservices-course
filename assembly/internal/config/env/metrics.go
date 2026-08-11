package env

import "github.com/caarlos0/env/v11"

type metricsEnvConfig struct {
	OTLPEndpoint string `env:"OTLP_ENDPOINT,required"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

func (cfg *metricsConfig) OTLPEndpoint() string {
	return cfg.raw.OTLPEndpoint
}
