// Package less provides comparison functions for various types.
package less

import "time"

// StringLess returns true if a is less than b.
func StringLess(a, b string) bool {
	return a < b
}

// IntLess returns true if a is less than b.
func IntLess(a, b int) bool {
	return a < b
}

// Float64Less returns true if a is less than b.
func Float64Less(a, b float64) bool {
	return a < b
}

// Float32Less returns true if a is less than b.
func Float32Less(a, b float32) bool {
	return a < b
}

// ByteLess returns true if a is less than b.
func ByteLess(a, b byte) bool {
	return a < b
}

// RuneLess returns true if a is less than b.
func RuneLess(a, b rune) bool {
	return a < b
}

// UintLess returns true if a is less than b.
func UintLess(a, b uint) bool {
	return a < b
}

// Uint8Less returns true if a is less than b.
func Uint8Less(a, b uint8) bool {
	return a < b
}

// Uint16Less returns true if a is less than b.
func Uint16Less(a, b uint16) bool {
	return a < b
}

// Uint32Less returns true if a is less than b.
func Uint32Less(a, b uint32) bool {
	return a < b
}

// Uint64Less returns true if a is less than b.
func Uint64Less(a, b uint64) bool {
	return a < b
}

// Int8Less returns true if a is less than b.
func Int8Less(a, b int8) bool {
	return a < b
}

// Int16Less returns true if a is less than b.
func Int16Less(a, b int16) bool {
	return a < b
}

// Int32Less returns true if a is less than b.
func Int32Less(a, b int32) bool {
	return a < b
}

// Int64Less returns true if a is less than b.
func Int64Less(a, b int64) bool {
	return a < b
}

// Complex64Less returns true if a is less than b, comparing real and imaginary parts.
func Complex64Less(a, b complex64) bool {
	return real(a) < real(b) || (real(a) == real(b) && imag(a) < imag(b))
}

// Complex128Less returns true if a is less than b, comparing real and imaginary parts.
func Complex128Less(a, b complex128) bool {
	return real(a) < real(b) || (real(a) == real(b) && imag(a) < imag(b))
}

// BoolLess returns true if a is false and b is true.
func BoolLess(a, b bool) bool {
	return !a && b
}

// TimeLess returns true if a is before b.
func TimeLess(a, b time.Time) bool {
	return a.Before(b)
}

// DurationLess returns true if a is less than b.
func DurationLess(a, b time.Duration) bool {
	return a < b
}
