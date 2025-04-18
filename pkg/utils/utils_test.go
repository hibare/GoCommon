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
