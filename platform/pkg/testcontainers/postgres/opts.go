package postgres

type Option func(*Config)

func WithNetworkName(network string) Option {
	return func(c *Config) {
		c.NetworkName = network
	}
}

func WithContainerName(containerName string) Option {
	return func(c *Config) {
		c.ContainerName = containerName
	}
}

// WithNetworkAliases задаёт дополнительные DNS-имена, по которым контейнер
// будет резолвиться внутри сети NetworkName — в отличие от ContainerName,
// не обязано совпадать с реальным именем контейнера в Docker.
func WithNetworkAliases(aliases ...string) Option {
	return func(c *Config) {
		c.NetworkAliases = aliases
	}
}

func WithImageName(image string) Option {
	return func(c *Config) {
		c.ImageName = image
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

func WithAuth(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
