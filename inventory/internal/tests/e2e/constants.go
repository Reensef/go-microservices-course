//go:build integration

package integration

const (
	// projectName - имя проекта для Docker-контейнеров и сети
	projectName = "inventory-service"

	// sightingsCollectionName - имя коллекции MongoDB для наблюдений НЛО
	partsCollectionName = "parts"

	// Конфигурация приложения внутри тестового окружения
	grpcHostValue     = "0.0.0.0"
	grpcPortValue     = "50051"
	loggerAsJSONValue = "true"

	// Конфигурация MongoDB внутри тестового окружения
	mongoImageNameValue = "mongo:7.0.5"
	mongoHostValue      = "inventory-mongo"
	mongoDatabaseValue  = "inventory"
	mongoAuthDBValue    = "admin"
	mongoUsernameValue  = "inventory_admin"
	mongoPasswordValue  = "inventory_secret" //nolint:gosec

	// Параметры контейнера iam — нужен, чтобы auth-interceptor inventory мог
	// провалидировать сессию через AuthService.Whoami
	iamAppName         = "iam-app"
	iamDockerfile      = "deploy/docker/iam/Dockerfile"
	iamGrpcHostValue   = "0.0.0.0"
	iamGrpcPortValue   = "50054"
	iamSessionTTLValue = "1h"
	iamMigrationsDir   = "./migrations"

	// Конфигурация Postgres для iam внутри тестового окружения
	iamPostgresImageNameValue = "postgres:16-alpine"
	iamPostgresHostValue      = "iam-postgres"
	iamPostgresDatabaseValue  = "iam"
	iamPostgresUsernameValue  = "iam"
	iamPostgresPasswordValue  = "iam_secret" //nolint:gosec

	// Конфигурация Redis для iam внутри тестового окружения
	iamRedisImageNameValue = "redis:7-alpine"
	iamRedisHostValue      = "iam-redis"
	iamRedisPasswordValue  = ""
	iamRedisDBValue        = "0"
)
