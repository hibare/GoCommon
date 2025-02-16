package db

import (
	"embed"
	"errors"
)

// DatabaseConfig holds the configuration for the database connection
type DatabaseConfig struct {
	DSN            string
	MigrationsPath string
	MigrationsFS   embed.FS
	DBType         Database
}

func (c DatabaseConfig) Validate() error {
	if c.DSN == "" {
		return errors.New("dsn is required")
	}

	if c.DBType == nil {
		return errors.New("dbtype is required")
	}

	return nil
}
