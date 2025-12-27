package db

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed test_data/migrations/*.sql
var testMigrationsPostgres embed.FS

func TestPostgresDatabase_Open(t *testing.T) {
	_, baseConfig, err := SetupMockPostgresDB(t)
	require.NoError(t, err)

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
				sqlDB, dErr := db.DB()
				if dErr == nil {
					_ = sqlDB.Ping()
				}
			}

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NotNil(t, db)

				// Additional connection verification for success case
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Ping())
			}
		})
	}
}

func TestPostgresDatabase_Migrate(t *testing.T) {
	_, baseConfig, err := SetupMockPostgresDB(t)
	require.NoError(t, err)

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
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify migration by checking if the test table exists
				var exists bool
				row := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tests')").Row()
				err := row.Scan(&exists)
				require.NoError(t, err)
				require.True(t, exists)
			}
		})
	}
}

func TestPostgresDatabase_RunSQLFromDirectory(t *testing.T) {
	_, baseConfig, err := SetupMockPostgresDB(t)
	require.NoError(t, err)

	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "successful_sql_execution",
			dir:     "test_data/sql_scripts",
			wantErr: false,
		},
		{
			name:    "invalid_directory",
			dir:     "invalid/dir",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgDB := &PostgresDatabase{}
			gormDB, err := pgDB.Open(baseConfig)
			require.NoError(t, err)

			// First run migrations to create the tests table
			err = pgDB.Migrate(gormDB, DatabaseConfig{
				DSN:            baseConfig.DSN,
				MigrationsFS:   testMigrationsPostgres,
				MigrationsPath: "test_data/migrations",
			})
			require.NoError(t, err)

			// Create DB wrapper
			db := &DB{
				DB:     gormDB,
				config: baseConfig,
			}

			err = db.RunSQLFromDirectory(tt.dir)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify view was created
				var exists bool
				row := gormDB.Raw("SELECT EXISTS (SELECT FROM information_schema.views WHERE table_name = 'test_view')").Row()
				err := row.Scan(&exists)
				require.NoError(t, err)
				require.True(t, exists)

				// Verify data was inserted
				var count int
				err = gormDB.Raw("SELECT COUNT(*) FROM tests WHERE name = 'test_from_sql'").Row().Scan(&count)
				require.NoError(t, err)
				require.Equal(t, 1, count)
			}
		})
	}
}
