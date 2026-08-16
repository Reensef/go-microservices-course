package testcontainers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

// Logger используется testcontainers-пакетами (mongo, postgres, redis) для логирования
// жизненного цикла контейнеров.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

// DefaultHostConfig возвращает модификатор HostConfig с автоматическим удалением
// контейнера после остановки — общий для всех testcontainers-пакетов.
func DefaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}

// GetContainerHostPort возвращает host и смаппленный порт запущенного контейнера
// для указанного внутреннего порта (в формате "5432/tcp").
func GetContainerHostPort(ctx context.Context, c testcontainers.Container, port string) (string, string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	mappedPort, err := c.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, mappedPort.Port(), nil
}
