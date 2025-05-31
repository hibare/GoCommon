package less

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringLess(t *testing.T) {
	assert.True(t, StringLess("a", "b"), "Expected 'a' to be less than 'b'")
	assert.False(t, StringLess("b", "a"), "Expected 'b' to not be less than 'a'")
}

func TestIntLess(t *testing.T) {
	assert.True(t, IntLess(1, 2), "Expected 1 to be less than 2")
	assert.False(t, IntLess(2, 1), "Expected 2 to not be less than 1")
}

func TestFloat64Less(t *testing.T) {
	assert.True(t, Float64Less(1.1, 2.2), "Expected 1.1 to be less than 2.2")
	assert.False(t, Float64Less(2.2, 1.1), "Expected 2.2 to not be less than 1.1")
}

func TestFloat32Less(t *testing.T) {
	assert.True(t, Float32Less(1.1, 2.2), "Expected 1.1 to be less than 2.2")
	assert.False(t, Float32Less(2.2, 1.1), "Expected 2.2 to not be less than 1.1")
}

func TestByteLess(t *testing.T) {
	assert.True(t, ByteLess(1, 2), "Expected 1 to be less than 2")
	assert.False(t, ByteLess(2, 1), "Expected 2 to not be less than 1")
}

func TestRuneLess(t *testing.T) {
	assert.True(t, RuneLess('a', 'b'), "Expected 'a' to be less than 'b'")
	assert.False(t, RuneLess('b', 'a'), "Expected 'b' to not be less than 'a'")
}

func TestUintLess(t *testing.T) {
	assert.True(t, UintLess(1, 2), "Expected 1 to be less than 2")
	assert.False(t, UintLess(2, 1), "Expected 2 to not be less than 1")
}

func TestUint8Less(t *testing.T) {
	assert.True(t, Uint8Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Uint8Less(2, 1), "Expected 2 to not be less than 1")
}

func TestUint16Less(t *testing.T) {
	assert.True(t, Uint16Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Uint16Less(2, 1), "Expected 2 to not be less than 1")
}

func TestUint32Less(t *testing.T) {
	assert.True(t, Uint32Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Uint32Less(2, 1), "Expected 2 to not be less than 1")
}

func TestUint64Less(t *testing.T) {
	assert.True(t, Uint64Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Uint64Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt8Less(t *testing.T) {
	assert.True(t, Int8Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Int8Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt16Less(t *testing.T) {
	assert.True(t, Int16Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Int16Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt32Less(t *testing.T) {
	assert.True(t, Int32Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Int32Less(2, 1), "Expected 2 to not be less than 1")
}

func TestInt64Less(t *testing.T) {
	assert.True(t, Int64Less(1, 2), "Expected 1 to be less than 2")
	assert.False(t, Int64Less(2, 1), "Expected 2 to not be less than 1")
}

func TestComplex64Less(t *testing.T) {
	assert.True(t, Complex64Less(1+1i, 2+2i), "Expected 1+1i to be less than 2+2i")
	assert.False(t, Complex64Less(2+2i, 1+1i), "Expected 2+2i to not be less than 1+1i")
}

func TestComplex128Less(t *testing.T) {
	assert.True(t, Complex128Less(1+1i, 2+2i), "Expected 1+1i to be less than 2+2i")
	assert.False(t, Complex128Less(2+2i, 1+1i), "Expected 2+2i to not be less than 1+1i")
}

func TestBoolLess(t *testing.T) {
	assert.True(t, BoolLess(false, true), "Expected false to be less than true")
	assert.False(t, BoolLess(true, false), "Expected true to not be less than false")
}

func TestTimeLess(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	assert.True(t, TimeLess(t1, t2), "Expected t1 to be less than t2")
	assert.False(t, TimeLess(t2, t1), "Expected t2 to not be less than t1")
}

func TestDurationLess(t *testing.T) {
	d1 := time.Duration(1)
	d2 := time.Duration(2)
	assert.True(t, DurationLess(d1, d2), "Expected d1 to be less than d2")
	assert.False(t, DurationLess(d2, d1), "Expected d2 to not be less than d1")
}
