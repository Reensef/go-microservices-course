//go:build integration

package environment

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/app"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/network"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/path"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/postgres"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/redis"
)

const (
	// Параметры для контейнеров
	iamAppName    = "iam-app"
	iamDockerfile = "deploy/docker/iam/Dockerfile"

	startupTimeout = 3 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network  *network.Network
	Postgres *postgres.Container
	Redis    *redis.Container
	App      *app.Container
}

// Setup — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами.
// Все настройки (порты, креды, имя образа) заданы в коде константами этого пакета,
// а не читаются из .env-файлов — тестовое окружение не зависит от deploy/compose.
func Setup(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Шаг 2: Запускаем контейнер с Postgres.
	// Алиас сети (postgresHostValue) — это DNS-имя, по которому приложение внутри
	// той же сети найдёт Postgres (передаётся ему через POSTGRES_HOST ниже).
	generatedPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithNetworkAliases(postgresHostValue),
		postgres.WithImageName(postgresImageNameValue),
		postgres.WithDatabase(postgresDatabaseValue),
		postgres.WithAuth(postgresUsernameValue, postgresPasswordValue),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanup(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер Postgres", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Postgres успешно запущен")

	// Шаг 3: Запускаем контейнер с Redis.
	generatedRedis, err := redis.NewContainer(ctx,
		redis.WithNetworkName(generatedNetwork.Name()),
		redis.WithContainerName(testcontainers.RedisContainerName),
		redis.WithNetworkAliases(redisHostValue),
		redis.WithImageName(redisImageNameValue),
		redis.WithPassword(redisPasswordValue),
		redis.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanup(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось запустить контейнер Redis", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Redis успешно запущен")

	// Шаг 4: Запускаем контейнер с приложением.
	// POSTGRES_PORT/REDIS_PORT — это порты *внутри* Docker-сети (testcontainers.PostgresPort,
	// testcontainers.RedisPort), а не внешние порты, проброшенные наружу для клиентов тестов.
	// Миграции БД приложение накатывает само при старте (см. iam/internal/app/app.go).
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		"GRPC_HOST":                        grpcHostValue,
		"GRPC_PORT":                        grpcPortValue,
		"LOGGER_LEVEL":                     loggerLevelValue,
		"LOGGER_AS_JSON":                   loggerAsJSONValue,
		"SESSION_TTL":                      sessionTTLValue,
		"MIGRATIONS_DIR":                   migrationsDirValue,
		testcontainers.PostgresHostKey:     postgresHostValue,
		testcontainers.PostgresPortKey:     testcontainers.PostgresPort,
		testcontainers.PostgresDatabaseKey: postgresDatabaseValue,
		testcontainers.PostgresUsernameKey: postgresUsernameValue,
		testcontainers.PostgresPasswordKey: postgresPasswordValue,
		testcontainers.RedisHostKey:        redisHostValue,
		testcontainers.RedisPortKey:        testcontainers.RedisPort,
		testcontainers.RedisPasswordKey:    redisPasswordValue,
		testcontainers.RedisDBKey:          redisDBValue,
	}

	waitStrategy := wait.ForListeningPort(nat.Port(grpcPortValue + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(iamAppName),
		app.WithPort(grpcPortValue),
		app.WithDockerfile(projectRoot, iamDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanup(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres, Redis: generatedRedis})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		Redis:    generatedRedis,
		App:      appContainer,
	}
}
