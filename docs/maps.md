# Maps Package Documentation

## Overview

The `maps` package provides utility functions for working with Go maps, including key/value extraction, sorting, file export, and conversion from `sync.Map`.

---

## Key Functions

- **MapContains(m, k) bool**: Checks if a key is present in a map.
- **MapKeys(m) []K**: Returns all keys in a map.
- **MapValues(m) []V**: Returns all values in a map.
- **MapSortByKeys(m, less) map[K]V**: Returns a new map sorted by keys.
- **MapSortByValues(m, less) map[K]V**: Returns a new map sorted by values.
- **Map2EnvFile(m, filePath) error**: Writes the map to a file as key=value pairs.
- **MapFromSyncMap(m \*sync.Map) map[K]V**: Converts a `sync.Map` to a regular map.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/maps"
)

m := map[string]int{"a": 1, "b": 2}
keys := maps.MapKeys(m)
```

---

## Notes

- Designed for use with Go generics.
- Useful for configuration, data export, and concurrent map handling.
