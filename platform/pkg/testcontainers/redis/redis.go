package redis

import (
	"context"
	"net"

	goredis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

type Container struct {
	container testcontainers.Container
	client    *goredis.Client
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startRedisContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate redis container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(cfg.Host, cfg.Port)

	client, err := connectRedisClient(ctx, addr, cfg.Password, cfg.Logger)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Redis container started", zap.String("addr", addr))
	success = true

	return &Container{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

func (c *Container) Client() *goredis.Client {
	return c.client
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.client.Close(); err != nil {
		c.cfg.Logger.Error(ctx, "failed to close redis client", zap.Error(err))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate redis container", zap.Error(err))
		return err
	}

	c.cfg.Logger.Info(ctx, "Redis container terminated")

	return nil
}
