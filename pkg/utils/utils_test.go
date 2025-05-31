package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostname(t *testing.T) {
	hostname := GetHostname()
	assert.NotEmpty(t, hostname)
}

func TestSyncMapLength(t *testing.T) {
	var m sync.Map

	// Test empty map
	if length := SyncMapLength(&m); length != 0 {
		t.Errorf("expected length 0, got %d", length)
	}

	// Test map with one element
	m.Store("key1", "value1")
	if length := SyncMapLength(&m); length != 1 {
		t.Errorf("expected length 1, got %d", length)
	}

	// Test map with multiple elements
	m.Store("key2", "value2")
	m.Store("key3", "value3")
	if length := SyncMapLength(&m); length != 3 {
		t.Errorf("expected length 3, got %d", length)
	}

	// Test map after deleting an element
	m.Delete("key2")
	if length := SyncMapLength(&m); length != 2 {
		t.Errorf("expected length 2, got %d", length)
	}
}

func TestToPtr(t *testing.T) {
	// int
	intVal := 42
	intPtr := ToPtr(intVal)
	assert.NotNil(t, intPtr)
	assert.Equal(t, intVal, *intPtr)

	// string
	strVal := "hello"
	strPtr := ToPtr(strVal)
	assert.NotNil(t, strPtr)
	assert.Equal(t, strVal, *strPtr)

	// bool
	boolVal := true
	boolPtr := ToPtr(boolVal)
	assert.NotNil(t, boolPtr)
	assert.Equal(t, boolVal, *boolPtr)

	// float64
	floatVal := 3.14
	floatPtr := ToPtr(floatVal)
	assert.NotNil(t, floatPtr)
	assert.Equal(t, floatVal, *floatPtr)

	// struct
	type sampleStruct struct {
		A int
		B string
	}
	structVal := sampleStruct{A: 1, B: "test"}
	structPtr := ToPtr(structVal)
	assert.NotNil(t, structPtr)
	assert.Equal(t, structVal, *structPtr)

	// array
	arrVal := [2]int{1, 2}
	arrPtr := ToPtr(arrVal)
	assert.NotNil(t, arrPtr)
	assert.Equal(t, arrVal, *arrPtr)

	// slice
	sliceVal := []string{"a", "b"}
	slicePtr := ToPtr(sliceVal)
	assert.NotNil(t, slicePtr)
	assert.Equal(t, sliceVal, *slicePtr)

	// map
	mapVal := map[string]int{"a": 1, "b": 2}
	mapPtr := ToPtr(mapVal)
	assert.NotNil(t, mapPtr)
	assert.Equal(t, mapVal, *mapPtr)

	// pointer
	origPtr := func() *int { v := 99; return &v }()
	ptrPtr := ToPtr(origPtr)
	assert.NotNil(t, ptrPtr)
	assert.Equal(t, *origPtr, **ptrPtr)

	// interface
	ifaceVal := interface{}("iface")
	ifacePtr := ToPtr(ifaceVal)
	assert.NotNil(t, ifacePtr)
	assert.Equal(t, ifaceVal, *ifacePtr)
}
