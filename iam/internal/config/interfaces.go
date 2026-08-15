package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OTLPEndpoint() string
}

type IamServiceConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
}

type SqlMigratorConfig interface {
	MigrationsDir() string
}

type RedisConfig interface {
	Address() string
	Password() string
	DB() int
}

type SessionConfig interface {
	TTL() time.Duration
}
