package redis

import (
	"context"

	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func connectRedisClient(ctx context.Context, addr, password string, logger Logger) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("failed to close redis client", zap.Error(closeErr))
		}
		return nil, errors.Errorf("failed to ping redis: %v", err)
	}

	return client, nil
}
