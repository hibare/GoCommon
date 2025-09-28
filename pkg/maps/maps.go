// Package maps provides utility functions for working with maps.
package maps

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

// Contains checks if a key is present in map m.
func Contains[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

// Keys returns all keys in a map.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values in a map.
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

// SortByKeys sorts a map by its keys and returns a new map with sorted keys.
func SortByKeys[K comparable, V any](m map[K]V, less func(a, b K) bool) map[K]V {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return less(keys[i], keys[j])
	})

	sortedMap := make(map[K]V, len(m))
	for _, k := range keys {
		sortedMap[k] = m[k]
	}

	return sortedMap
}

// SortByValues sorts a map by its values and returns a new map with sorted values.
func SortByValues[K comparable, V any](m map[K]V, less func(a, b V) bool) map[K]V {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return less(m[keys[i]], m[keys[j]])
	})

	sortedMap := make(map[K]V, len(m))
	for _, k := range keys {
		sortedMap[k] = m[k]
	}

	return sortedMap
}

// ToEnvFile writes the map to a file where each line is key=value.
func ToEnvFile[K comparable, V any](m map[K]V, filePath string) error {
	// Open or create the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// Iterate over the map and write each key-value pair to the file
	for key, value := range m {
		if _, err := fmt.Fprintf(file, "%v=%v\n", key, value); err != nil {
			return err
		}
	}

	return nil
}

// FromSyncMap converts a sync.Map to a regular map.
func FromSyncMap[K comparable, V any](m *sync.Map) map[K]V {
	result := make(map[K]V)
	m.Range(func(key, value any) bool {
		k, ok := key.(K)
		if !ok {
			return true
		}
		v, ok := value.(V)
		if !ok {
			return true
		}
		result[k] = v
		return true
	})
	return result
}
