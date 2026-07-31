#!/bin/sh
set -e

# Выполняется один раз при первом старте контейнера (когда $PGDATA пуст).
# Создаёт базы, перечисленные через запятую в POSTGRES_MULTIPLE_DATABASES,
# по одной на каждый микросервис, использующий этот общий инстанс Postgres.

if [ -z "$POSTGRES_MULTIPLE_DATABASES" ]; then
  exit 0
fi

OLD_IFS=$IFS
IFS=','
for db in $POSTGRES_MULTIPLE_DATABASES; do
  echo "Creating database '$db' (if not exists) for user '$POSTGRES_USER'"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    SELECT 'CREATE DATABASE "$db" OWNER "$POSTGRES_USER"'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$db')\gexec
EOSQL
done
IFS=$OLD_IFS
