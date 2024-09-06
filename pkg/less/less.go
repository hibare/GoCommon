package less

import "time"

func StringLess(a, b string) bool {
	return a < b
}

func IntLess(a, b int) bool {
	return a < b
}

func Float64Less(a, b float64) bool {
	return a < b
}

func Float32Less(a, b float32) bool {
	return a < b
}

func ByteLess(a, b byte) bool {
	return a < b
}

func RuneLess(a, b rune) bool {
	return a < b
}

func UintLess(a, b uint) bool {
	return a < b
}

func Uint8Less(a, b uint8) bool {
	return a < b
}

func Uint16Less(a, b uint16) bool {
	return a < b
}

func Uint32Less(a, b uint32) bool {
	return a < b
}

func Uint64Less(a, b uint64) bool {
	return a < b
}

func Int8Less(a, b int8) bool {
	return a < b
}

func Int16Less(a, b int16) bool {
	return a < b
}

func Int32Less(a, b int32) bool {
	return a < b
}

func Int64Less(a, b int64) bool {
	return a < b
}

func Complex64Less(a, b complex64) bool {
	return real(a) < real(b) || (real(a) == real(b) && imag(a) < imag(b))
}

func Complex128Less(a, b complex128) bool {
	return real(a) < real(b) || (real(a) == real(b) && imag(a) < imag(b))
}

func BoolLess(a, b bool) bool {
	return !a && b
}

func TimeLess(a, b time.Time) bool {
	return a.Before(b)
}

func DurationLess(a, b time.Duration) bool {
	return a < b
}
