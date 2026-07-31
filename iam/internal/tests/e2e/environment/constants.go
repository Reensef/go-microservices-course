//go:build integration

package environment

const (
	// projectName - имя проекта для Docker-контейнеров и сети
	projectName = "iam-service"

	// usersTableName - имя таблицы Postgres с пользователями
	usersTableName = "users"

	// Конфигурация приложения внутри тестового окружения
	grpcHostValue      = "0.0.0.0"
	grpcPortValue      = "50054"
	loggerAsJSONValue  = "true"
	loggerLevelValue   = "debug"
	sessionTTLValue    = "1h"
	migrationsDirValue = "./migrations"

	// Конфигурация Postgres внутри тестового окружения
	postgresImageNameValue = "postgres:16-alpine"
	postgresHostValue      = "iam-postgres"
	postgresDatabaseValue  = "iam"
	postgresUsernameValue  = "iam"
	postgresPasswordValue  = "iam_secret" //nolint:gosec

	// Конфигурация Redis внутри тестового окружения
	redisImageNameValue = "redis:7-alpine"
	redisHostValue      = "iam-redis"
	redisPasswordValue  = ""
	redisDBValue        = "0"
)
