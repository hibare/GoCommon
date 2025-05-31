# Env Package Documentation

## Overview

The `env` package provides utilities for working with environment variables in Go applications. It supports loading `.env` files, retrieving environment variables with type conversion, and extracting prefixed variables as maps.

---

## Key Functions

- **Load() error**: Loads an optional `.env` file using [godotenv](https://github.com/joho/godotenv).
- **MustString(key, fallback) string**: Returns the value of an environment variable or a fallback.
- **MustBool(key, fallback) bool**: Returns a boolean from an environment variable or fallback.
- **MustInt(key, fallback) int**: Returns an integer from an environment variable or fallback.
- **MustDuration(key, fallback) time.Duration**: Returns a duration from an environment variable or fallback.
- **MustStringSlice(key, fallback) []string**: Returns a string slice from a comma-separated environment variable or fallback.
- **GetPrefixed(prefix) map[string]string**: Returns a map of environment variables with the given prefix.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/env"
)

_ = env.Load()
port := env.MustInt("PORT", 8080)
```

---

## Notes

- Provides type-safe access to environment variables with sensible fallbacks.
- Useful for configuration in 12-factor applications.
