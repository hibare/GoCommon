package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hibare/GoCommon/v2/pkg/constants"
	"github.com/hibare/GoCommon/v2/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	configRootLinux   = constants.ConfigRootLinux
	configRootWindows = constants.ConfigRootWindows
	configRootDarwin  = constants.ConfigRootDarwin
)

type OSInterface interface {
	GetGOOS() string
}

type ActualOS struct{}

func (ao ActualOS) GetGOOS() string {
	return runtime.GOOS
}

type BaseConfig struct {
	ProgramIdentifier   string
	ConfigRootDir       string
	ConfigFilePath      string
	OS                  OSInterface
	ConfigFileName      string
	ConfigFileExtension string
}

func (bc *BaseConfig) SetConfigRootDir() {
	var configRootDir string

	if os.Getenv("IS_LOCAL") == "true" {
		cwd, err := os.Getwd()
		if err != nil {
			return
		}
		configRootDir = filepath.Join(cwd, "/etc/")
	} else {
		switch bc.OS.GetGOOS() {
		case "linux":
			configRootDir = configRootLinux
		case "windows":
			configRootDir = configRootWindows
		case "darwin":
			configRootDir = configRootDarwin
		default:
			bc.ConfigRootDir = ""
			return
		}
	}

	bc.ConfigRootDir = filepath.Join(configRootDir, strings.ToLower(bc.ProgramIdentifier))
}

func (bc *BaseConfig) SetConfigFilePath() {
	bc.SetConfigRootDir()
	bc.ConfigFilePath = filepath.Join(bc.ConfigRootDir, fmt.Sprintf("%s.%s", constants.ConfigFileName, constants.ConfigFileExtension))
}

func (bc *BaseConfig) EnsureConfigRootDir() error {
	if bc.ConfigRootDir == "" {
		bc.SetConfigRootDir()
	}

	if info, err := os.Stat(bc.ConfigRootDir); os.IsNotExist(err) {
		if err := os.MkdirAll(bc.ConfigRootDir, 0755); err != nil {
			return errors.ErrCreatingDir
		}
	} else if !info.IsDir() {
		return errors.ErrNotDir
	}

	return nil
}

func (bc *BaseConfig) EnsureConfigFile() error {
	if bc.ConfigFilePath == "" {
		bc.SetConfigFilePath()
	}
	err := bc.EnsureConfigRootDir()

	if err != nil {
		return err
	}

	if info, err := os.Stat(bc.ConfigFilePath); os.IsNotExist(err) {
		file, err := os.Create(bc.ConfigFilePath)
		if err != nil {
			return err
		}
		file.Close()
	} else if info.IsDir() {
		return errors.ErrNotFile
	}

	return nil
}

func (bc *BaseConfig) Init() error {
	bc.ConfigFileName = constants.ConfigFileName
	bc.ConfigFileExtension = constants.ConfigFileExtension
	bc.SetConfigRootDir()
	bc.SetConfigFilePath()
	if err := bc.EnsureConfigRootDir(); err != nil {
		return err
	}

	if err := bc.EnsureConfigFile(); err != nil {
		return err
	}
	return nil
}

func (bc *BaseConfig) WriteYAMLConfig(current interface{}) error {
	v := viper.New()
	v.SetConfigType(bc.ConfigFileExtension)
	v.SetConfigName(bc.ConfigFileName)
	v.AddConfigPath(bc.ConfigRootDir)

	yamlData, err := yaml.Marshal(current)
	if err != nil {
		return err
	}

	if err = v.ReadConfig(bytes.NewBuffer(yamlData)); err != nil {
		return err
	}

	if err = v.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func (bc *BaseConfig) ReadYAMLConfig(current interface{}) (interface{}, error) {
	v := viper.New()
	v.SetConfigType(bc.ConfigFileExtension)
	v.SetConfigName(bc.ConfigFileName)
	v.AddConfigPath(bc.ConfigRootDir)

	err := v.ReadInConfig()
	if err != nil {
		return current, err
	}

	if err := v.Unmarshal(&current); err != nil {
		return current, err
	}

	return current, nil
}

func (bc *BaseConfig) CleanConfig() error {
	if info, err := os.Stat(bc.ConfigRootDir); err != nil && os.IsExist(err) {
		return err
	} else if !info.IsDir() {
		return errors.ErrNotDir
	} else {
		if err := os.RemoveAll(bc.ConfigRootDir); err != nil {
			return err
		}
	}
	return nil
}
