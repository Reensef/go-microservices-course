package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" // nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)

// Postgres constants
const (
	// Postgres container constants
	PostgresContainerName = "postgres"
	PostgresPort          = "5432"

	// Postgres environment variables
	PostgresHostKey     = "POSTGRES_HOST"
	PostgresPortKey     = "POSTGRES_PORT"
	PostgresDatabaseKey = "POSTGRES_DATABASE"
	PostgresUsernameKey = "POSTGRES_USER"
	PostgresPasswordKey = "POSTGRES_PASSWORD" // nolint:gosec
)

// Redis constants
const (
	// Redis container constants
	RedisContainerName = "redis"
	RedisPort          = "6379"

	// Redis environment variables
	RedisHostKey     = "REDIS_HOST"
	RedisPortKey     = "REDIS_PORT"
	RedisPasswordKey = "REDIS_PASSWORD" // nolint:gosec
	RedisDBKey       = "REDIS_DB"
)
