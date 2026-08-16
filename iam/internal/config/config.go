package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Reensef/go-microservices-course/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	IamService  IamServiceConfig
	Postgres    PostgresConfig
	SqlMigrator SqlMigratorConfig
	Redis       RedisConfig
	Session     SessionConfig
	Tracing     TracingConfig
}

func Load(envFile string) error {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	iamServiceCfg, err := env.NewIamServiceConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	sqlMigratorCfg, err := env.NewSqlMigratorConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		IamService:  iamServiceCfg,
		Postgres:    postgresCfg,
		SqlMigrator: sqlMigratorCfg,
		Redis:       redisCfg,
		Session:     sessionCfg,
		Tracing:     tracingCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
