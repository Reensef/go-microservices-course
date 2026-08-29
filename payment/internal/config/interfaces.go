package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
}

type PaymentServiceConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceVersion() string
}

type ServiceConfig interface {
	Name() string
	Environment() string
}
