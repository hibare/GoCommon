package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/hibare/GoCommon/v2/pkg/constants"
	"github.com/hibare/GoCommon/v2/pkg/errors"
	"github.com/stretchr/testify/require"
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
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.EnsureConfigRootDir()
	require.NoError(t, err)
	require.DirExists(t, bc.ConfigRootDir)
	_ = os.Remove(bc.ConfigRootDir)
}

func TestEnsureConfigRootDirFail(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	f, err := os.Create(bc.ConfigRootDir)
	require.NoError(t, err)
	_ = f.Close()
	err = bc.EnsureConfigRootDir()
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrNotDir)
	_ = os.Remove(bc.ConfigRootDir)
}

func TestEnsureConfigFile(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.EnsureConfigFile()
	require.NoError(t, err)
	require.FileExists(t, bc.ConfigFilePath)
	_ = os.RemoveAll(bc.ConfigRootDir)
}

func TestEnsureConfigFileFail(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	bc.SetConfigFilePath()
	err := os.MkdirAll(bc.ConfigFilePath, 0755)
	require.NoError(t, err)
	err = bc.EnsureConfigFile()
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrNotFile)
	_ = os.RemoveAll(bc.ConfigRootDir)
}

func TestInitBC(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.Init()
	require.NoError(t, err)
	require.DirExists(t, bc.ConfigRootDir)
	require.FileExists(t, bc.ConfigFilePath)
	_ = os.RemoveAll(bc.ConfigRootDir)
}

func TestInitBCFailRootDir(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	f, err := os.Create(bc.ConfigRootDir)
	require.NoError(t, err)
	_ = f.Close()
	err = bc.Init()
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrNotDir)
	_ = os.RemoveAll(bc.ConfigRootDir)
}

func TestInitBCFailConfigFile(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	bc.SetConfigRootDir()
	bc.SetConfigFilePath()
	err := os.MkdirAll(bc.ConfigFilePath, 0755)
	require.NoError(t, err)
	err = bc.Init()
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrNotFile)
	_ = os.RemoveAll(bc.ConfigRootDir)
}

func TestWriteReadYAMLConfig(t *testing.T) {
	configRootLinux = t.TempDir()
	bc := BaseConfig{
		ProgramIdentifier: "myapp",
		OS:                &MockOS{MockGOOS: "linux"},
	}
	err := bc.Init()
	require.NoError(t, err)
	require.DirExists(t, bc.ConfigRootDir)
	require.FileExists(t, bc.ConfigFilePath)

	type config struct {
		ProgramIdentifier string `yaml:"programIdentifier" mapstructure:"programIdentifier"`
	}
	writeConfig := config{
		ProgramIdentifier: "myapp",
	}
	err = bc.WriteYAMLConfig(writeConfig)
	require.NoError(t, err)

	var readConfig *config
	rConfig, err := bc.ReadYAMLConfig(readConfig)
	readConfig, ok := rConfig.(*config)
	require.True(t, ok)
	require.NoError(t, err)
	require.Equal(t, writeConfig, *readConfig)

	_ = os.RemoveAll(bc.ConfigRootDir)
}
