package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePG "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDatabase implements the Database interface for PostgreSQL.
type PostgresDatabase struct{}

// Open opens a database connection.
func (p *PostgresDatabase) Open(config DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
}

// Migrate migrates the database.
func (p *PostgresDatabase) Migrate(db *gorm.DB, config DatabaseConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := migratePG.WithInstance(sqlDB, &migratePG.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	source, err := iofs.New(config.MigrationsFS, config.MigrationsPath)
	if err != nil {
		return fmt.Errorf("could not create source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}
