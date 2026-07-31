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
)
