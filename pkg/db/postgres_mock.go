// Package db provides utilities for working with databases.
package db

import (
	"fmt"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

const (
	// PGTestUser is the test user for the postgres database.
	PGTestUser = "test_user"

	// PGTestPass is the test password for the postgres database.
	PGTestPass = "test_pass"

	// PGTestDB is the test database for the postgres database.
	PGTestDB = "test_db"

	// PGVersion is the version of the postgres database.
	PGVersion = "15"
)

// SetupMockPostgresDB sets up a mock postgres database.
func SetupMockPostgresDB() (*gnomock.Container, DatabaseConfig, error) {
	p := postgres.Preset(
		postgres.WithUser(PGTestUser, PGTestPass),
		postgres.WithDatabase(PGTestDB),
		postgres.WithVersion(PGVersion),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		return nil, DatabaseConfig{}, err
	}

	return container, DatabaseConfig{
		DSN: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			container.Host,
			container.DefaultPort(),
			PGTestUser,
			PGTestPass,
			PGTestDB,
		),
	}, nil
}

// UnsetMockPostgresDB unsets the mock postgres database.
func UnsetMockPostgresDB(container *gnomock.Container) error {
	return gnomock.Stop(container)
}
