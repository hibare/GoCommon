# Slice Package Documentation

## Overview

The `slice` package provides utilities for working with Go slices, including set-like operations such as unique, diff, intersect, and union. These functions are generic and work with any comparable type.

---

## Key Functions

- **Unique(list []T) []T**: Returns a slice with unique elements.
- **Diff(a, b []T) []T**: Returns the elements in `a` that are not in `b`.
- **Intersect(a, b []T) []T**: Returns the intersection of two slices.
- **Union(a, b []T) []T**: Returns the union of two slices.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/slice"
)

unique := slice.Unique([]int{1, 2, 2, 3}) // [1, 2, 3]
diff := slice.Diff([]int{1, 2, 3}, []int{2, 3}) // [1]
inter := slice.Intersect([]int{1, 2, 3}, []int{2, 3, 4}) // [2, 3]
union := slice.Union([]int{1, 2}, []int{2, 3}) // [1, 2, 3]
```

---

## Notes

- Designed for use with Go generics.
- Useful for set operations and data processing.
