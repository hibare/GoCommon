# Utils Package Documentation

## Overview

The `utils` package provides miscellaneous utility functions for Go applications, including hostname retrieval, sync.Map length calculation, and a generic pointer helper.

---

## Key Functions

- **GetHostname() string**: Returns the hostname of the machine.
- **SyncMapLength(m \*sync.Map) int**: Returns the number of elements in a `sync.Map`.
- **ToPtr(v T) \*T**: Returns a pointer to the provided value (generic helper).

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/utils"
)

hostname := utils.GetHostname()
ptr := utils.ToPtr(42)
```

---

## Notes

- Designed for general-purpose utility needs in Go projects.
