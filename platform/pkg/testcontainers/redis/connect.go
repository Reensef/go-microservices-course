package redis

import (
	"context"

	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
)

func connectRedisClient(ctx context.Context, addr, password string) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, errors.Errorf("failed to ping redis: %v", err)
	}

	return client, nil
}
