package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OTLPEndpoint() string
}

type PaymentServiceConfig interface {
	Address() string
}

type IAMClientConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
