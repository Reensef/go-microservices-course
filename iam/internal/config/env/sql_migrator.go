package env

import (
	"github.com/caarlos0/env/v11"
)

type sqlMigratorEnvConfig struct {
	MigrationsDir string `env:"MIGRATIONS_DIR,required"`
}

type sqlMigratorConfig struct {
	raw sqlMigratorEnvConfig
}

func NewSqlMigratorConfig() (*sqlMigratorConfig, error) {
	var raw sqlMigratorEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &sqlMigratorConfig{raw: raw}, nil
}

func (cfg *sqlMigratorConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
