package db

import (
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"
)

var (
	ErrUnsupportedDriver = errors.New("unsupported database driver")
)

// Database interface defines the methods for database operations
type Database interface {
	Open(config DatabaseConfig) (*gorm.DB, error)
	Migrate(db *gorm.DB, config DatabaseConfig) error
}

type DB struct {
	DB     *gorm.DB
	config DatabaseConfig
}

var (
	instance *DB
	mu       sync.Mutex
)

func (d *DB) Migrate() error {
	return d.config.DBType.Migrate(d.DB, d.config)
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	instance = nil
	return sqlDB.Close()
}

// NewClient returns a singleton instance of the database connection
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
		sqlDB.Close()
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
