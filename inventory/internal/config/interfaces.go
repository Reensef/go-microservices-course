package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryServiceConfig interface {
	Address() string
}

type IAMClientConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
