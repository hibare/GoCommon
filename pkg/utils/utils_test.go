package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHostname(t *testing.T) {
	hostname := GetHostname()
	require.NotEmpty(t, hostname)
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
	require.NotNil(t, intPtr)
	require.Equal(t, intVal, *intPtr)

	// string
	strVal := "hello"
	strPtr := ToPtr(strVal)
	require.NotNil(t, strPtr)
	require.Equal(t, strVal, *strPtr)

	// bool
	boolVal := true
	boolPtr := ToPtr(boolVal)
	require.NotNil(t, boolPtr)
	require.Equal(t, boolVal, *boolPtr)

	// float64
	floatVal := 3.14
	floatPtr := ToPtr(floatVal)
	require.NotNil(t, floatPtr)
	require.InEpsilon(t, floatVal, *floatPtr, 1e-9)

	// struct
	type sampleStruct struct {
		A int
		B string
	}
	structVal := sampleStruct{A: 1, B: "test"}
	structPtr := ToPtr(structVal)
	require.NotNil(t, structPtr)
	require.Equal(t, structVal, *structPtr)

	// array
	arrVal := [2]int{1, 2}
	arrPtr := ToPtr(arrVal)
	require.NotNil(t, arrPtr)
	require.Equal(t, arrVal, *arrPtr)

	// slice
	sliceVal := []string{"a", "b"}
	slicePtr := ToPtr(sliceVal)
	require.NotNil(t, slicePtr)
	require.Equal(t, sliceVal, *slicePtr)

	// map
	mapVal := map[string]int{"a": 1, "b": 2}
	mapPtr := ToPtr(mapVal)
	require.NotNil(t, mapPtr)
	require.Equal(t, mapVal, *mapPtr)

	// pointer
	origPtr := func() *int { v := 99; return &v }()
	ptrPtr := ToPtr(origPtr)
	require.NotNil(t, ptrPtr)
	require.Equal(t, *origPtr, **ptrPtr)

	// interface
	ifaceVal := interface{}("iface")
	ifacePtr := ToPtr(ifaceVal)
	require.NotNil(t, ifacePtr)
	require.Equal(t, ifaceVal, *ifacePtr)
}
