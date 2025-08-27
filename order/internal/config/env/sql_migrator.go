package env

import (
	"github.com/caarlos0/env/v11"
)

type sqlMigratorEnvConfig struct {
	MigrationsDir string `env:"MIGRATIONS_DIR,required"`
}

type SqlMigratorConfig struct {
	raw sqlMigratorEnvConfig
}

func NewSqlMigratorConfig() (*SqlMigratorConfig, error) {
	var raw sqlMigratorEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &SqlMigratorConfig{raw: raw}, nil
}

func (cfg *SqlMigratorConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
