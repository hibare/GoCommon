# Less Package Documentation

## Overview

The `less` package provides comparison functions for various Go types. These functions are useful for sorting and ordering operations, especially when working with generics or custom sort logic.

---

## Key Functions

- **StringLess(a, b string) bool**: Returns true if `a < b`.
- **IntLess(a, b int) bool**: Returns true if `a < b`.
- **Float64Less(a, b float64) bool**: Returns true if `a < b`.
- **Float32Less(a, b float32) bool**: Returns true if `a < b`.
- **ByteLess(a, b byte) bool**: Returns true if `a < b`.
- **RuneLess(a, b rune) bool**: Returns true if `a < b`.
- **UintLess(a, b uint) bool**: Returns true if `a < b`.
- **Uint8Less(a, b uint8) bool**: Returns true if `a < b`.
- **Uint16Less(a, b uint16) bool**: Returns true if `a < b`.
- **Uint32Less(a, b uint32) bool**: Returns true if `a < b`.
- **Uint64Less(a, b uint64) bool**: Returns true if `a < b`.
- **Int8Less(a, b int8) bool**: Returns true if `a < b`.
- **Int16Less(a, b int16) bool**: Returns true if `a < b`.
- **Int32Less(a, b int32) bool**: Returns true if `a < b`.
- **Int64Less(a, b int64) bool**: Returns true if `a < b`.
- **Complex64Less(a, b complex64) bool**: Compares real and imaginary parts.
- **Complex128Less(a, b complex128) bool**: Compares real and imaginary parts.
- **BoolLess(a, b bool) bool**: Returns true if `a` is false and `b` is true.
- **TimeLess(a, b time.Time) bool**: Returns true if `a` is before `b`.
- **DurationLess(a, b time.Duration) bool**: Returns true if `a < b`.

---

## Example Usage

```go
import (
    "sort"
    "github.com/hibare/GoCommon/v2/pkg/less"
)

strings := []string{"b", "a", "c"}
sort.Slice(strings, func(i, j int) bool { return less.StringLess(strings[i], strings[j]) })
```

---

## Notes

- Designed for use with Go's `sort` package and generics.
