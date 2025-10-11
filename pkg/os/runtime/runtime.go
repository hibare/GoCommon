package runtime

import (
	"fmt"
	"path/filepath"
	actualRuntime "runtime"
)

const (
	// ConfigRootLinux is the config root directory for Linux.
	ConfigRootLinux = "/etc/"

	// ConfigRootWindows is the config root directory for Windows.
	ConfigRootWindows = "C:\\ProgramData\\"

	// ConfigRootDarwin is the config root directory for Darwin.
	ConfigRootDarwin = "/Library/Application Support/"

	// ConfigFileName is the config file name.
	ConfigFileName = "config"

	// ConfigFileExtension is the config file extension.
	ConfigFileExtension = "yaml"
)

// RuntimeIface is the interface for the runtime.
type RuntimeIface interface {
	GetGOOS() string
	GetPlatform() string
	GetConfigDir() string
	GetConfigFilePath() string
}

// Runtime is the actual runtime.
type Runtime struct{}

// GetGOOS returns the GOOS.
func (Runtime) GetGOOS() string {
	return actualRuntime.GOOS
}

// GetPlatform returns the platform in the format os/arch.
func (Runtime) GetPlatform() string {
	return actualRuntime.GOOS + "/" + actualRuntime.GOARCH
}

// GetConfigDir returns the default config directory based on the OS.
func (r Runtime) GetConfigDir() string {
	switch r.GetGOOS() {
	case "windows":
		return ConfigRootWindows
	case "darwin":
		return ConfigRootDarwin
	default:
		return ConfigRootLinux
	}
}

// GetConfigFilePath returns the full path to the config file.
func (r Runtime) GetConfigFilePath() string {
	fileName := fmt.Sprintf("%s.%s", ConfigFileName, ConfigFileExtension)
	return filepath.Join(r.GetConfigDir(), fileName)
}

func newRuntime() RuntimeIface {
	return Runtime{}
}

// New is the actual runtime instance.
var New = newRuntime
