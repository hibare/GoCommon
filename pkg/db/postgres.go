package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePG "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDatabase implements the Database interface for PostgreSQL
type PostgresDatabase struct{}

func (p *PostgresDatabase) Open(config DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
}

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

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}
