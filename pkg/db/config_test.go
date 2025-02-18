package db

import (
	"testing"
)

func TestDatabaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  DatabaseConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: DatabaseConfig{
				DSN:    "valid_dsn",
				DBType: &MockDatabase{},
			},
			wantErr: false,
		},
		{
			name: "missing DSN",
			config: DatabaseConfig{
				DBType: &MockDatabase{},
			},
			wantErr: true,
		},
		{
			name: "missing DBType",
			config: DatabaseConfig{
				DSN: "valid_dsn",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
