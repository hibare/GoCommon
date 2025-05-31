# Datetime Package Documentation

## Overview

The `datetime` package provides utilities for working with date and time strings, including sorting date-time values in a specific format.

---

## Key Functions

- **SortDateTimes(dt []string) []string**: Sorts a slice of date-time strings (using the default layout from `constants`). Returns the sorted slice in descending order (most recent first).

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/datetime"
)

dates := []string{"20230721053000", "20230720053000", "20230722053000"}
sorted := datetime.SortDateTimes(dates)
// sorted = ["20230722053000", "20230721053000", "20230720053000"]
```

---

## Notes

- Expects date-time strings in the format defined by `constants.DefaultDateTimeLayout` (e.g., `20060102150405`).
