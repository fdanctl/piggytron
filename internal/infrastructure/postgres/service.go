// Package postgres implements the repository and query-service interfaces on
// top of PostgreSQL with raw SQL. Repositories accept a DBTX (*sql.DB or
// *sql.Tx) so they can join application-level transactions; query services
// implement the read-model interfaces declared in internal/query.
package postgres

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var (
	ErrMissingMigrationsPath = errors.New("MIGRATIONS_PATH env missing")
	ErrMissingDatabaseURL    = errors.New("DATABASE_URL env missing")
)

//go:embed migrations/*.sql
var migrations embed.FS

// Connect opens a db and runs migrations
func Connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
	if err != nil {
		return nil, fmt.Errorf("migrate new: %s", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
