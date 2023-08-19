package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/hibare/GoCommon/v2/pkg/constants"
	"github.com/hibare/GoCommon/v2/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type MockOS struct {
	MockGOOS string
}

func (m MockOS) GetGOOS() string {
	return m.MockGOOS
}

func TestGetConfigRootDir(t *testing.T) {
	testCases := []struct {
		name              string
		programIdentifier string
		OS                OSInterface
		expectedResult    string
	}{
		{
			name:              "Linux OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "linux"},
			expectedResult:    "/etc/myapp",
		},
		{
			name:              "Windows OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "windows"},
			expectedResult:    "C:\\ProgramData\\/myapp",
		},
		{
			name:              "Darwin OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "darwin"},
			expectedResult:    "/Library/Application Support/myapp",
		},
		{
			name:              "Unknown OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "unknown"},
			expectedResult:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bc := BaseConfig{
				ProgramIdentifier: tc.programIdentifier,
				OS:                tc.OS,
			}
			bc.SetConfigRootDir()
			if bc.ConfigRootDir != tc.expectedResult {
				t.Errorf("Expected %s, but got %s", tc.expectedResult, bc.ConfigRootDir)
			}
		})
	}
}

func TestGetConfigFilePath(t *testing.T) {
	testCases := []struct {
		name              string
		programIdentifier string
		OS                OSInterface
		expectedResult    string
	}{
		{
			name:              "Linux OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "linux"},
			expectedResult:    fmt.Sprintf("/etc/myapp/%s.%s", constants.ConfigFileName, constants.ConfigFileExtension),
		},
		{
			name:              "Windows OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "windows"},
			expectedResult:    fmt.Sprintf("C:\\ProgramData\\/myapp/%s.%s", constants.ConfigFileName, constants.ConfigFileExtension),
		},
		{
			name:              "Darwin OS",
			programIdentifier: "myapp",
			OS:                &MockOS{MockGOOS: "darwin"},
			expectedResult:    fmt.Sprintf("/Library/Application Support/myapp/%s.%s", constants.ConfigFileName, constants.ConfigFileExtension),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bc := BaseConfig{
				ProgramIdentifier: tc.programIdentifier,
				OS:                tc.OS,
			}
			bc.SetConfigFilePath()
			if bc.ConfigFilePath != tc.expectedResult {
				t.Errorf("Expected %s, but got %s", tc.expectedResult, bc.ConfigFilePath)
			}
		})
	}
}

func TestEnsureConfigRootDir(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.EnsureConfigRootDir()
	assert.NoError(t, err)
	assert.DirExists(t, bc.ConfigRootDir)
	os.Remove(bc.ConfigRootDir)
}

func TestEnsureConfigRootDirFail(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	f, err := os.Create(bc.ConfigRootDir)
	assert.NoError(t, err)
	f.Close()
	err = bc.EnsureConfigRootDir()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrNotDir)
	os.Remove(bc.ConfigRootDir)
}

func TestEnsureConfigFile(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.EnsureConfigFile()
	assert.NoError(t, err)
	assert.FileExists(t, bc.ConfigFilePath)
	os.RemoveAll(bc.ConfigRootDir)
}

func TestEnsureConfigFileFail(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	bc.SetConfigFilePath()
	err := os.MkdirAll(bc.ConfigFilePath, 0755)
	assert.NoError(t, err)
	err = bc.EnsureConfigFile()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrNotFile)
	os.RemoveAll(bc.ConfigRootDir)
}

func TestInitBC(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.Init()
	assert.NoError(t, err)
	assert.DirExists(t, bc.ConfigRootDir)
	assert.FileExists(t, bc.ConfigFilePath)
	os.RemoveAll(bc.ConfigRootDir)
}

func TestInitBCFailRootDir(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	f, err := os.Create(bc.ConfigRootDir)
	assert.NoError(t, err)
	f.Close()
	err = bc.Init()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrNotDir)
	os.RemoveAll(bc.ConfigRootDir)
}

func TestInitBCFailConfigFile(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	bc.SetConfigFilePath()
	err := os.MkdirAll(bc.ConfigFilePath, 0755)
	assert.NoError(t, err)
	err = bc.Init()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrNotFile)
	os.RemoveAll(bc.ConfigRootDir)
}

func TestWriteReadYAMLConfig(t *testing.T) {
	configRootLinux = os.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.Init()
	assert.NoError(t, err)
	assert.DirExists(t, bc.ConfigRootDir)
	assert.FileExists(t, bc.ConfigFilePath)

	type config struct {
		ProgramIdentifier string `yaml:"programIdentifier" mapstructure:"programIdentifier"`
	}
	writeConfig := config{
		ProgramIdentifier: "myapp",
	}
	err = bc.WriteYAMLConfig(writeConfig)
	assert.NoError(t, err)

	var readConfig *config
	rConfig, err := bc.ReadYAMLConfig(readConfig)
	readConfig = rConfig.(*config)
	assert.NoError(t, err)
	assert.Equal(t, writeConfig, *readConfig)

	os.RemoveAll(bc.ConfigRootDir)
}
