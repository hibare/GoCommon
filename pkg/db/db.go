// Package db provides database abstraction and utilities.
package db

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// RunSQLFromDirectory executes all .sql files found in the specified directory.
// Files are executed in alphabetical order.
func (d *DB) RunSQLFromDirectory(dir string) error {
	var files []string

	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Sort files alphabetically for consistent execution order
	sort.Strings(files)

	// Execute each SQL file
	for _, file := range files {
		if err := d.executeSQLFile(file); err != nil {
			return err
		}
	}

	return nil
}

// RunSQLFromFS executes all .sql files found in the embedded filesystem directory.
// Files are executed in alphabetical order.
func (d *DB) RunSQLFromFS(fsys fs.FS, dir string) error {
	var files []string

	err := fs.WalkDir(fsys, dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Sort files alphabetically for consistent execution order
	sort.Strings(files)

	// Execute each SQL file
	for _, file := range files {
		if err := d.executeSQLFileFS(fsys, file); err != nil {
			return err
		}
	}

	return nil
}

// executeSQLFile reads and executes the SQL content from a filesystem file.
func (d *DB) executeSQLFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	sql := strings.TrimSpace(string(content))
	if sql == "" {
		return nil
	}

	return d.DB.Exec(sql).Error
}

// executeSQLFileFS reads and executes the SQL content from an embedded FS file.
func (d *DB) executeSQLFileFS(fsys fs.FS, filePath string) error {
	content, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return err
	}

	sql := strings.TrimSpace(string(content))
	if sql == "" {
		return nil
	}

	return d.DB.Exec(sql).Error
}
