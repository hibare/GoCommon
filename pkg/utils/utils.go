// Package utils provides utilities for working with utils.
package utils

import (
	"os"
	"sync"
)

// GetHostname returns the hostname of the machine.
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// SyncMapLength returns the number of elements in a sync.Map.
func SyncMapLength(m *sync.Map) int {
	count := 0
	m.Range(func(_, _ any) bool {
		count++
		return true // Continue iteration
	})
	return count
}

// ToPtr returns a pointer to the provided value v.
//
// This generic utility is useful for constructing pointer values
// for literals or values in a concise and type-safe manner.
// Example usage:
//
//	ptr := ToPtr(42) // ptr is of type *int
func ToPtr[T any](v T) *T {
	return &v
}
