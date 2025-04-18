package utils

import (
	"os"
	"sync"
)

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// SyncMapLength returns the number of elements in a sync.Map
func SyncMapLength(m *sync.Map) int {
	count := 0
	m.Range(func(_, _ any) bool {
		count++
		return true // Continue iteration
	})
	return count
}
