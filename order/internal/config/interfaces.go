package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderServiceConfig interface {
	Address() string
}

type PaymentClientConfig interface {
	Address() string
}

type InventoryClientConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type PostgresConfig interface {
	URI() string
}

type SqlMigratorConfig interface {
	MigrationsDir() string
}
