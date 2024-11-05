package structs

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Field1 string
	Field2 int
	Field3 bool
}

func TestStructCopy(t *testing.T) {
	src := &TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}
	dst := &TestStruct{}

	err := StructCopy(src, dst)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("expected dst to be %v, got %v", src, dst)
	}
}

func TestStructCompare(t *testing.T) {
	a := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}
	b := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}
	c := TestStruct{
		Field1: "different",
		Field2: 456,
		Field3: false,
	}

	equal, err := StructCompare(a, b)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !equal {
		t.Fatalf("expected structs to be equal")
	}

	equal, err = StructCompare(a, c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equal {
		t.Fatalf("expected structs to be different")
	}
}

func TestStructToMap(t *testing.T) {
	s := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}

	m, err := StructToMap(s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := map[string]interface{}{
		"Field1": "test",
		"Field2": 123,
		"Field3": true,
	}

	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("expected map to be %v, got %v", expected, m)
	}
}

func TestStructContainsField(t *testing.T) {
	s := TestStruct{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}

	contains, err := StructContainsField(s, "Field1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !contains {
		t.Fatalf("expected struct to contain field 'Field1'")
	}

	contains, err = StructContainsField(s, "NonExistentField")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if contains {
		t.Fatalf("expected struct to not contain field 'NonExistentField'")
	}
}
