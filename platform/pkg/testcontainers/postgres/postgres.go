package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

type Container struct {
	container testcontainers.Container
	pool      *pgxpool.Pool
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startPostgresContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	uri := buildPostgresURI(cfg)

	pool, err := connectPostgresPool(ctx, uri)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Postgres container started", zap.String("uri", uri))
	success = true

	return &Container{
		container: container,
		pool:      pool,
		cfg:       cfg,
	}, nil
}

func (c *Container) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	c.pool.Close()

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
		return err
	}

	c.cfg.Logger.Info(ctx, "Postgres container terminated")

	return nil
}
