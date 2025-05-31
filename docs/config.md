# Config Package Documentation

## Overview

The `config` package provides utilities for loading, managing, and validating configuration files in Go applications. It supports YAML configuration, OS-specific config root directories, and ensures config files and directories exist as needed.

---

## Key Types and Functions

- **BaseConfig**: Struct for managing configuration file paths, root directories, and OS abstraction.
- **OSInterface / ActualOS**: Interface and implementation for retrieving the current OS (for testability).
- **SetConfigRootDir()**: Sets the root directory for config files based on OS.
- **SetConfigFilePath()**: Sets the full path to the config file.
- **EnsureConfigRootDir()**: Ensures the config root directory exists.
- **EnsureConfigFile()**: Ensures the config file exists.
- **Init()**: Initializes the config, ensuring all paths and files are set up.
- **WriteYAMLConfig(current interface{})**: Writes a struct as YAML to the config file.
- **ReadYAMLConfig(current interface{})**: Reads YAML config into a struct.
- **CleanConfig()**: Removes the config directory and all contents.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/config"
)

bc := config.BaseConfig{ProgramIdentifier: "myapp", OS: config.ActualOS{}}
err := bc.Init()
if err != nil {
    panic(err)
}
// Write config
err = bc.WriteYAMLConfig(myConfigStruct)
// Read config
var cfg MyConfig
_, err = bc.ReadYAMLConfig(&cfg)
```

---

## Notes

- Supports Linux, Windows, and macOS config root conventions.
- Uses [viper](https://github.com/spf13/viper) and [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) for config management.
- Designed for testability and extensibility.
