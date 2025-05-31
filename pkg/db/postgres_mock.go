package db

import (
	"fmt"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

const (
	PGTestUser = "test_user"
	PGTestPass = "test_pass"
	PGTestDB   = "test_db"
	PGVersion  = "15"
)

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

func UnsetMockPostgresDB(container *gnomock.Container) error {
	return gnomock.Stop(container)
}
