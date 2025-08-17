package utils

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type SqlMigrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewSqlMigrator(db *sql.DB, migrationsDir string) *SqlMigrator {
	return &SqlMigrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *SqlMigrator) Up() error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
