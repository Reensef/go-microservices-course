package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentServiceConfig interface {
	Address() string
}
