//go:build integration

package integration

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
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/mongo"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/network"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/path"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/postgres"
	"github.com/Reensef/go-microservices-course/platform/pkg/testcontainers/redis"
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container

	IamPostgres *postgres.Container
	IamRedis    *redis.Container
	IamApp      *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами.
// Все настройки (порты, креды, имя образа) заданы в коде константами этого пакета,
// а не читаются из .env-файлов — тестовое окружение не зависит от deploy/compose.
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Шаг 2: Запускаем контейнер с MongoDB
	// Алиас сети (mongoHostValue) — это DNS-имя, по которому приложение внутри
	// той же сети найдёт Mongo (передаётся ему через MONGO_HOST ниже).
	// Само имя контейнера при этом остаётся отдельным (testcontainers.MongoContainerName),
	// чтобы не конфликтовать с контейнером Mongo из локального docker-compose разработчика.
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithNetworkAliases(mongoHostValue),
		mongo.WithImageName(mongoImageNameValue),
		mongo.WithDatabase(mongoDatabaseValue),
		mongo.WithAuth(mongoUsernameValue, mongoPasswordValue),
		mongo.WithAuthDB(mongoAuthDBValue),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	// Шаг 3: Запускаем контейнеры iam (Postgres + Redis + приложение) — auth-interceptor
	// inventory ходит в AuthService.Whoami, поэтому валидному запросу к inventory
	// нужен настоящий iam с действующей сессией.
	generatedIamPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithNetworkAliases(iamPostgresHostValue),
		postgres.WithImageName(iamPostgresImageNameValue),
		postgres.WithDatabase(iamPostgresDatabaseValue),
		postgres.WithAuth(iamPostgresUsernameValue, iamPostgresPasswordValue),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер Postgres для iam", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Postgres для iam успешно запущен")

	generatedIamRedis, err := redis.NewContainer(ctx,
		redis.WithNetworkName(generatedNetwork.Name()),
		redis.WithNetworkAliases(iamRedisHostValue),
		redis.WithImageName(iamRedisImageNameValue),
		redis.WithPassword(iamRedisPasswordValue),
		redis.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network: generatedNetwork, Mongo: generatedMongo, IamPostgres: generatedIamPostgres,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер Redis для iam", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Redis для iam успешно запущен")

	projectRoot := path.GetProjectRoot()

	iamAppEnv := map[string]string{
		"GRPC_HOST":                        iamGrpcHostValue,
		"GRPC_PORT":                        iamGrpcPortValue,
		"LOGGER_LEVEL":                     loggerLevelValue,
		"LOGGER_AS_JSON":                   loggerAsJSONValue,
		"SESSION_TTL":                      iamSessionTTLValue,
		"MIGRATIONS_DIR":                   iamMigrationsDir,
		testcontainers.PostgresHostKey:     iamPostgresHostValue,
		testcontainers.PostgresPortKey:     testcontainers.PostgresPort,
		testcontainers.PostgresDatabaseKey: iamPostgresDatabaseValue,
		testcontainers.PostgresUsernameKey: iamPostgresUsernameValue,
		testcontainers.PostgresPasswordKey: iamPostgresPasswordValue,
		testcontainers.RedisHostKey:        iamRedisHostValue,
		testcontainers.RedisPortKey:        testcontainers.RedisPort,
		testcontainers.RedisPasswordKey:    iamRedisPasswordValue,
		testcontainers.RedisDBKey:          iamRedisDBValue,
	}

	iamWaitStrategy := wait.ForListeningPort(nat.Port(iamGrpcPortValue + "/tcp")).
		WithStartupTimeout(startupTimeout)

	iamAppContainer, err := app.NewContainer(ctx,
		app.WithName(iamAppName),
		app.WithPort(iamGrpcPortValue),
		app.WithDockerfile(projectRoot, iamDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(iamAppEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(iamWaitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network: generatedNetwork, Mongo: generatedMongo,
			IamPostgres: generatedIamPostgres, IamRedis: generatedIamRedis,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения iam", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения iam успешно запущен")

	// Шаг 4: Запускаем контейнер с приложением.
	// MONGO_PORT — это порт Mongo *внутри* Docker-сети (testcontainers.MongoPort),
	// а не внешний порт, который testcontainers пробрасывает наружу для клиента тестов.
	// IAM_CLIENT_HOST/PORT указывают на контейнер iam-app в этой же сети — auth-interceptor
	// inventory использует их, чтобы независимо валидировать сессию через AuthService.Whoami.
	appEnv := map[string]string{
		"GRPC_HOST":                     grpcHostValue,
		"GRPC_PORT":                     grpcPortValue,
		"LOGGER_LEVEL":                  loggerLevelValue,
		"LOGGER_AS_JSON":                loggerAsJSONValue,
		"IAM_CLIENT_HOST":               iamAppName,
		"IAM_CLIENT_PORT":               iamGrpcPortValue,
		testcontainers.MongoHostKey:     mongoHostValue,
		testcontainers.MongoPortKey:     testcontainers.MongoPort,
		testcontainers.MongoDatabaseKey: mongoDatabaseValue,
		testcontainers.MongoAuthDBKey:   mongoAuthDBValue,
		testcontainers.MongoUsernameKey: mongoUsernameValue,
		testcontainers.MongoPasswordKey: mongoPasswordValue,
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPortValue + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPortValue),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network: generatedNetwork, Mongo: generatedMongo,
			IamPostgres: generatedIamPostgres, IamRedis: generatedIamRedis, IamApp: iamAppContainer,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:     generatedNetwork,
		Mongo:       generatedMongo,
		IamPostgres: generatedIamPostgres,
		IamRedis:    generatedIamRedis,
		IamApp:      iamAppContainer,
		App:         appContainer,
	}
}
