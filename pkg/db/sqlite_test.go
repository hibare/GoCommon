package db

import (
	"embed"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

//go:embed test_data/migrations/*.sql
var testMigrationsSqlite embed.FS

func TestSQLiteDatabase_Open(t *testing.T) {
	tests := []struct {
		name    string
		config  DatabaseConfig
		wantErr bool
	}{
		{
			name: "success_memory_db",
			config: DatabaseConfig{
				DSN: "file::memory:?cache=shared",
			},
			wantErr: false,
		},
		{
			name: "success_file_db",
			config: DatabaseConfig{
				DSN: "test.db",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &SQLiteDatabase{}
			gormDB, err := db.Open(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, gormDB)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gormDB)

				// Test connection
				sqlDB, err := gormDB.DB()
				require.NoError(t, err)
				assert.NoError(t, sqlDB.Ping())

				// Cleanup if file-based database
				if tt.config.DSN != "file::memory:?cache=shared" {
					err := sqlDB.Close()
					require.NoError(t, err)
					err = os.Remove(tt.config.DSN)
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestSQLiteDatabase_Migrate(t *testing.T) {
	tests := []struct {
		name    string
		config  DatabaseConfig
		setup   func(*testing.T) *gorm.DB
		wantErr bool
	}{
		{
			name: "successful_migration",
			config: DatabaseConfig{
				DSN:            "file::memory:?cache=shared",
				MigrationsFS:   testMigrationsSqlite,
				MigrationsPath: "test_data/migrations",
			},
			setup: func(t *testing.T) *gorm.DB {
				db := &SQLiteDatabase{}
				gormDB, err := db.Open(DatabaseConfig{DSN: "file::memory:?cache=shared"})
				require.NoError(t, err)
				return gormDB
			},
			wantErr: false,
		},
		{
			name: "invalid_migrations_path",
			config: DatabaseConfig{
				DSN:            "file::memory:?cache=shared",
				MigrationsFS:   testMigrationsSqlite,
				MigrationsPath: "invalid/path",
			},
			setup: func(t *testing.T) *gorm.DB {
				db := &SQLiteDatabase{}
				gormDB, err := db.Open(DatabaseConfig{DSN: "file::memory:?cache=shared"})
				require.NoError(t, err)
				return gormDB
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &SQLiteDatabase{}
			gormDB := tt.setup(t)
			err := db.Migrate(gormDB, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSQLiteDatabase_Integration(t *testing.T) {
	// Create test table migration
	const createTableSQL = `
    CREATE TABLE IF NOT EXISTS tests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	db := &SQLiteDatabase{}
	config := DatabaseConfig{
		DSN: "file::memory:?cache=shared",
	}

	// Open database
	gormDB, err := db.Open(config)
	require.NoError(t, err)

	// Execute test migration
	sqlDB, err := gormDB.DB()
	require.NoError(t, err)
	_, err = sqlDB.Exec(createTableSQL)
	require.NoError(t, err)

	// Test data insertion and retrieval
	type Test struct {
		ID        uint      `gorm:"primarykey"`
		Name      string    `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	// Create
	test := &Test{Name: "test_record"}
	result := gormDB.Create(test)
	assert.NoError(t, result.Error)
	assert.NotZero(t, test.ID)

	// Read
	var retrieved Test
	result = gormDB.First(&retrieved, test.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, test.Name, retrieved.Name)
}
