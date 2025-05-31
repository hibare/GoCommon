package less

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStringLess(t *testing.T) {
	require.True(t, StringLess("a", "b"), "Expected 'a' to be less than 'b'")
	require.False(t, StringLess("b", "a"), "Expected 'b' to not be less than 'a'")
}

func TestIntLess(t *testing.T) {
	require.True(t, IntLess(1, 2), "Expected 1 to be less than 2")
	require.False(t, IntLess(2, 1), "Expected 2 to not be less than 1")
}

func TestFloat64Less(t *testing.T) {
	require.True(t, Float64Less(1.1, 2.2), "Expected 1.1 to be less than 2.2")
	require.False(t, Float64Less(2.2, 1.1), "Expected 2.2 to not be less than 1.1")
}

func TestFloat32Less(t *testing.T) {
	require.True(t, Float32Less(1.1, 2.2), "Expected 1.1 to be less than 2.2")
	require.False(t, Float32Less(2.2, 1.1), "Expected 2.2 to not be less than 1.1")
}

func TestByteLess(t *testing.T) {
	require.True(t, ByteLess(1, 2), "Expected 1 to be less than 2")
	require.False(t, ByteLess(2, 1), "Expected 2 to not be less than 1")
}

func TestRuneLess(t *testing.T) {
	require.True(t, RuneLess('a', 'b'), "Expected 'a' to be less than 'b'")
	require.False(t, RuneLess('b', 'a'), "Expected 'b' to not be less than 'a'")
}

func TestUintLess(t *testing.T) {
	require.True(t, UintLess(1, 2), "Expected 1 to be less than 2")
	require.False(t, UintLess(2, 1), "Expected 2 to not be less than 1")
}

func TestUint8Less(t *testing.T) {
	require.True(t, Uint8Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Uint8Less(2, 1), "Expected 2 to not be less than 1")
}

func TestUint16Less(t *testing.T) {
	require.True(t, Uint16Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Uint16Less(2, 1), "Expected 2 to not be less than 1")
}

func TestUint32Less(t *testing.T) {
	require.True(t, Uint32Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Uint32Less(2, 1), "Expected 2 to not be less than 1")
}

func TestUint64Less(t *testing.T) {
	require.True(t, Uint64Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Uint64Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt8Less(t *testing.T) {
	require.True(t, Int8Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Int8Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt16Less(t *testing.T) {
	require.True(t, Int16Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Int16Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt32Less(t *testing.T) {
	require.True(t, Int32Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Int32Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt64Less(t *testing.T) {
	require.True(t, Int64Less(1, 2), "Expected 1 to be less than 2")
	require.False(t, Int64Less(2, 1), "Expected 2 to not be less than 1")
}

func TestComplex64Less(t *testing.T) {
	require.True(t, Complex64Less(1+1i, 2+2i), "Expected 1+1i to be less than 2+2i")
	require.False(t, Complex64Less(2+2i, 1+1i), "Expected 2+2i to not be less than 1+1i")
}

func TestComplex128Less(t *testing.T) {
	require.True(t, Complex128Less(1+1i, 2+2i), "Expected 1+1i to be less than 2+2i")
	require.False(t, Complex128Less(2+2i, 1+1i), "Expected 2+2i to not be less than 1+1i")
}

func TestBoolLess(t *testing.T) {
	require.True(t, BoolLess(false, true), "Expected false to be less than true")
	require.False(t, BoolLess(true, false), "Expected true to not be less than false")
}

func TestTimeLess(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	require.True(t, TimeLess(t1, t2), "Expected t1 to be less than t2")
	require.False(t, TimeLess(t2, t1), "Expected t2 to not be less than t1")
}

func TestDurationLess(t *testing.T) {
	d1 := time.Duration(1)
	d2 := time.Duration(2)
	require.True(t, DurationLess(d1, d2), "Expected d1 to be less than d2")
	require.False(t, DurationLess(d2, d1), "Expected d2 to not be less than d1")
}
