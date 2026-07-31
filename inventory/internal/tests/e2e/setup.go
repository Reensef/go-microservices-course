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

	// Шаг 3: Запускаем контейнер с приложением.
	// MONGO_PORT — это порт Mongo *внутри* Docker-сети (testcontainers.MongoPort),
	// а не внешний порт, который testcontainers пробрасывает наружу для клиента тестов.
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		"GRPC_HOST":                     grpcHostValue,
		"GRPC_PORT":                     grpcPortValue,
		"LOGGER_LEVEL":                  loggerLevelValue,
		"LOGGER_AS_JSON":                loggerAsJSONValue,
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
