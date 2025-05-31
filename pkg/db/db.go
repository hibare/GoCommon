// Package db provides database abstraction and utilities.
package db

import (
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"
)

// ErrUnsupportedDriver is returned when an unsupported database driver is used.
var ErrUnsupportedDriver = errors.New("unsupported database driver")

// Database interface defines the methods for database operations.
type Database interface {
	Open(config DatabaseConfig) (*gorm.DB, error)
	Migrate(db *gorm.DB, config DatabaseConfig) error
}

// DB wraps a gorm.DB and its configuration.
type DB struct {
	DB     *gorm.DB
	config DatabaseConfig
}

var (
	instance *DB
	mu       sync.Mutex
)

// Migrate runs the migration for the database.
func (d *DB) Migrate() error {
	return d.config.DBType.Migrate(d.DB, d.config)
}

// Close closes the database connection.
func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	instance = nil
	return nil
}

// NewClient returns a singleton instance of the database connection.
func NewClient(ctx context.Context, config DatabaseConfig) (*DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		// Check if the existing connection is still alive
		sqlDB, err := instance.DB.DB()
		if err == nil && sqlDB.Ping() == nil {
			// Connection is alive, return the existing instance
			return instance, nil
		}
		// Connection is not alive, close it and create a new one
		_ = sqlDB.Close()
		instance = nil
	}

	db, err := config.DBType.Open(config)
	if err != nil {
		return nil, err
	}

	db = db.WithContext(ctx)
	instance = &DB{
		DB:     db,
		config: config,
	}

	return instance, nil
}
