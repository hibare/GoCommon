package less

import (
	"testing"
	"time"
)

func TestStringLess(t *testing.T) {
	if !StringLess("a", "b") {
		t.Errorf("Expected 'a' to be less than 'b'")
	}
	if StringLess("b", "a") {
		t.Errorf("Expected 'b' to not be less than 'a'")
	}
}

func TestIntLess(t *testing.T) {
	if !IntLess(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if IntLess(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestFloat64Less(t *testing.T) {
	if !Float64Less(1.1, 2.2) {
		t.Errorf("Expected 1.1 to be less than 2.2")
	}
	if Float64Less(2.2, 1.1) {
		t.Errorf("Expected 2.2 to not be less than 1.1")
	}
}

func TestFloat32Less(t *testing.T) {
	if !Float32Less(1.1, 2.2) {
		t.Errorf("Expected 1.1 to be less than 2.2")
	}
	if Float32Less(2.2, 1.1) {
		t.Errorf("Expected 2.2 to not be less than 1.1")
	}
}

func TestByteLess(t *testing.T) {
	if !ByteLess(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if ByteLess(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestRuneLess(t *testing.T) {
	if !RuneLess('a', 'b') {
		t.Errorf("Expected 'a' to be less than 'b'")
	}
	if RuneLess('b', 'a') {
		t.Errorf("Expected 'b' to not be less than 'a'")
	}
}

func TestUintLess(t *testing.T) {
	if !UintLess(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if UintLess(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestUint8Less(t *testing.T) {
	if !Uint8Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Uint8Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestUint16Less(t *testing.T) {
	if !Uint16Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Uint16Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestUint32Less(t *testing.T) {
	if !Uint32Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Uint32Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestUint64Less(t *testing.T) {
	if !Uint64Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Uint64Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestInt8Less(t *testing.T) {
	if !Int8Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Int8Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestInt16Less(t *testing.T) {
	if !Int16Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Int16Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestInt32Less(t *testing.T) {
	if !Int32Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Int32Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestInt64Less(t *testing.T) {
	if !Int64Less(1, 2) {
		t.Errorf("Expected 1 to be less than 2")
	}
	if Int64Less(2, 1) {
		t.Errorf("Expected 2 to not be less than 1")
	}
}

func TestComplex64Less(t *testing.T) {
	if !Complex64Less(1+1i, 2+2i) {
		t.Errorf("Expected 1+1i to be less than 2+2i")
	}
	if Complex64Less(2+2i, 1+1i) {
		t.Errorf("Expected 2+2i to not be less than 1+1i")
	}
}

func TestComplex128Less(t *testing.T) {
	if !Complex128Less(1+1i, 2+2i) {
		t.Errorf("Expected 1+1i to be less than 2+2i")
	}
	if Complex128Less(2+2i, 1+1i) {
		t.Errorf("Expected 2+2i to not be less than 1+1i")
	}
}

func TestBoolLess(t *testing.T) {
	if !BoolLess(false, true) {
		t.Errorf("Expected false to be less than true")
	}
	if BoolLess(true, false) {
		t.Errorf("Expected true to not be less than false")
	}
}

func TestTimeLess(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	if !TimeLess(t1, t2) {
		t.Errorf("Expected t1 to be less than t2")
	}
	if TimeLess(t2, t1) {
		t.Errorf("Expected t2 to not be less than t1")
	}
}

func TestDurationLess(t *testing.T) {
	d1 := time.Duration(1)
	d2 := time.Duration(2)
	if !DurationLess(d1, d2) {
		t.Errorf("Expected d1 to be less than d2")
	}
	if DurationLess(d2, d1) {
		t.Errorf("Expected d2 to not be less than d1")
	}
}
