package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	return gormDB
}

func TestNewClient(t *testing.T) {
	// Setup
	gormDB := setupTestDB(t)
	mockDB := new(MockDatabase)
	config := DatabaseConfig{
		DSN:    ":memory:",
		DBType: mockDB,
	}

	// Set up expectations with mock.Anything for gormDB comparison
	mockDB.On("Open", config).Return(gormDB, nil).Once()
	mockDB.On("Migrate", mock.Anything, config).Return(nil).Once()

	// Test first client creation
	client, err := NewClient(t.Context(), config)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Instead of comparing entire GORM instances, verify the connection works
	sqlDB1, err := gormDB.DB()
	require.NoError(t, err)
	sqlDB2, err := client.DB.DB()
	require.NoError(t, err)
	assert.NoError(t, sqlDB1.Ping())
	assert.NoError(t, sqlDB2.Ping())

	err = client.Migrate()
	assert.NoError(t, err)

	// Test singleton behavior
	mockDB.On("Open", config).Return(gormDB, nil).Maybe()
	client2, err := NewClient(t.Context(), config)
	require.NoError(t, err)
	assert.Equal(t, client, client2) // This works because it's the same instance

	// Close DB
	err = client.Close()
	require.NoError(t, err)
	// Verify all expectations were met
	mockDB.AssertExpectations(t)
}

func TestNewClient_OpenError(t *testing.T) {
	mockDB := new(MockDatabase)
	config := DatabaseConfig{
		DSN:    ":memory:",
		DBType: mockDB,
	}

	openErr := errors.New("failed to open database")
	mockDB.On("Open", config).Return(nil, openErr)

	client, err := NewClient(t.Context(), config)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, openErr, err)

	mockDB.AssertExpectations(t)
}

func TestNewClient_MigrateError(t *testing.T) {
	mockDB := new(MockDatabase)
	config := DatabaseConfig{
		DSN:    ":memory:",
		DBType: mockDB,
	}

	// Create a new GORM DB instance for testing
	gormDB, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{})
	assert.NoError(t, err)

	mockDB.On("Open", config).Return(gormDB, nil)
	migrateErr := errors.New("failed to migrate database")
	mockDB.On("Migrate", mock.Anything, config).Return(migrateErr)

	client, err := NewClient(t.Context(), config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	err = client.Migrate()
	assert.Error(t, err)
	assert.Equal(t, migrateErr, err)

	mockDB.AssertExpectations(t)
}

func TestDB_Close(t *testing.T) {
	// Create a new GORM DB instance for testing
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	client := &DB{DB: gormDB}
	err = client.Close()
	assert.NoError(t, err)
}
