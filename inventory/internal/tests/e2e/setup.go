//go:build integration

package integration

import (
	"context"
	"maps"
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
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context, envVars map[string]string) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, envVars, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, envVars, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, envVars, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, envVars, testcontainers.MongoDatabaseKey)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, envVars, grpcPortKey)

	// Шаг 2: Запускаем контейнер с MongoDB
	// Алиас сети совпадает с MONGO_HOST из .env, поэтому приложению не нужно
	// объяснять контейнеру, где искать Mongo, — .env работает как есть, без подмен.
	// Само имя контейнера при этом остаётся отдельным (testcontainers.MongoContainerName),
	// чтобы не конфликтовать с контейнером Mongo из локального docker-compose разработчика.
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithNetworkAliases(envVars[testcontainers.MongoHostKey]),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	appEnv := maps.Clone(envVars)

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения из .env с логированием при отсутствии
func getEnvWithLogging(ctx context.Context, envVars map[string]string, key string) string {
	value := envVars[key]
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
