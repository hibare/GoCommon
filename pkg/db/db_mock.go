package db

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Open(config DatabaseConfig) (*gorm.DB, error) {
	args := m.Called(config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	gormDB, ok := args.Get(0).(*gorm.DB)
	if !ok {
		return nil, args.Error(1)
	}
	return gormDB, args.Error(1)
}

func (m *MockDatabase) Migrate(gormDB *gorm.DB, config DatabaseConfig) error {
	args := m.Called(gormDB, config)
	return args.Error(0)
}
