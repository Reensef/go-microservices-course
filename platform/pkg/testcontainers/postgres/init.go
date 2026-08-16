package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	tc "github.com/Reensef/go-microservices-course/platform/pkg/testcontainers"
)

const (
	postgresPort           = "5432"
	postgresStartupTimeout = 1 * time.Minute

	postgresEnvUserKey     = "POSTGRES_USER"
	postgresEnvPasswordKey = "POSTGRES_PASSWORD" //nolint:gosec
	postgresEnvDatabaseKey = "POSTGRES_DB"
)

func startPostgresContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		NetworkAliases: map[string][]string{
			cfg.NetworkName: cfg.NetworkAliases,
		},
		Env: map[string]string{
			postgresEnvUserKey:     cfg.Username,
			postgresEnvPasswordKey: cfg.Password,
			postgresEnvDatabaseKey: cfg.Database,
		},
		// Официальный образ postgres дважды пишет это сообщение в лог: первый раз
		// после временного запуска для initdb, второй раз — когда сервер готов
		// принимать внешние соединения. Ждём именно второе вхождение.
		WaitingFor:         wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(postgresStartupTimeout),
		HostConfigModifier: tc.DefaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start postgres container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	return tc.GetContainerHostPort(ctx, container, postgresPort+"/tcp")
}

func buildPostgresURI(cfg *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}
