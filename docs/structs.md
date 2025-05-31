# Structs Package Documentation

## Overview

The `structs` package provides utilities for working with Go structs, including copying, comparing, converting to maps, and checking for fields. These functions use reflection and are useful for generic struct manipulation.

---

## Key Functions

- **StructCopy(src, dst interface{}) error**: Copies the contents of one struct to another (both must be pointers to structs).
- **StructCompare(a, b interface{}) (bool, error)**: Compares two structs for equality.
- **StructToMap(s interface{}) (map[string]interface{}, error)**: Converts a struct to a map.
- **StructContainsField(s interface{}, fieldName string) (bool, error)**: Checks if a struct contains a field with the given name.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/structs"
)

type MyStruct struct {
    Field1 string
    Field2 int
}
a := MyStruct{"foo", 1}
b := MyStruct{"foo", 1}
equal, _ := structs.StructCompare(a, b) // true
```

---

## Notes

- Deprecated: `CopyStruct` is an alias for `StructCopy` and will be removed in the future.
- Useful for generic struct manipulation and testing.
