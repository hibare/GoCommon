package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migrateSQLite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLiteDatabase implements the Database interface for SQLite.
type SQLiteDatabase struct{}

// Open opens a database connection.
func (s *SQLiteDatabase) Open(config DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(config.DSN), &gorm.Config{})
}

// Migrate migrates the database.
func (s *SQLiteDatabase) Migrate(db *gorm.DB, config DatabaseConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := migrateSQLite.WithInstance(sqlDB, &migrateSQLite.Config{})
	if err != nil {
		return fmt.Errorf("could not create sqlite driver: %w", err)
	}

	source, err := iofs.New(config.MigrationsFS, config.MigrationsPath)
	if err != nil {
		return fmt.Errorf("could not create source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}
