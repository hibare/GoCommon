package db

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed test_data/migrations/*.sql
var testMigrationsPostgres embed.FS

func TestPostgresDatabase(t *testing.T) {
	container, baseConfig, err := SetupMockPostgresDB()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = UnsetMockPostgresDB(container)
	})

	t.Run("Open", func(t *testing.T) {
		tests := []struct {
			name    string
			config  DatabaseConfig
			wantErr bool
		}{
			{
				name:    "success",
				config:  baseConfig,
				wantErr: false,
			},
			{
				name: "connection_error",
				config: DatabaseConfig{
					DSN: "invalid-dsn",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pgDB := &PostgresDatabase{}
				db, err := pgDB.Open(tt.config)

				// For connection errors, we need to attempt to use the connection
				if err == nil {
					sqlDB, err := db.DB()
					if err == nil {
						err = sqlDB.Ping()
					}
				}

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, db)

					// Additional connection verification for success case
					sqlDB, err := db.DB()
					require.NoError(t, err)
					assert.NoError(t, sqlDB.Ping())
				}
			})
		}
	})

	t.Run("Migrate", func(t *testing.T) {
		tests := []struct {
			name    string
			config  DatabaseConfig
			wantErr bool
		}{
			{
				name: "successful_migration",
				config: DatabaseConfig{
					DSN:            baseConfig.DSN,
					MigrationsFS:   testMigrationsPostgres,
					MigrationsPath: "test_data/migrations",
				},
				wantErr: false,
			},
			{
				name: "invalid_migrations_path",
				config: DatabaseConfig{
					DSN:            baseConfig.DSN,
					MigrationsFS:   testMigrationsPostgres,
					MigrationsPath: "invalid/path",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pgDB := &PostgresDatabase{}
				db, err := pgDB.Open(tt.config)
				require.NoError(t, err)

				err = pgDB.Migrate(db, tt.config)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					// Verify migration by checking if the test table exists
					var exists bool
					row := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tests')").Row()
					err := row.Scan(&exists)
					require.NoError(t, err)
					assert.True(t, exists)
				}
			})
		}
	})
}
