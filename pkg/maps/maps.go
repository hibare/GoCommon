package maps

import (
	"fmt"
	"os"
	"sort"
)

// MapContains checks if a key is present in map
func MapContains[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

// MapKeys returns all keys in a map
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// MapValues return all values in a map
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

// MapSortByKeys sorts a map by its keys and returns a slice of key-value pairs
func MapSortByKeys[K comparable, V any](m map[K]V, less func(a, b K) bool) map[K]V {
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

func MapSortByValues[K comparable, V any](m map[K]V, less func(a, b V) bool) map[K]V {
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

// Map2EnvFile writes the map to a file where each line is key=value
func Map2EnvFile[K comparable, V any](m map[K]V, filePath string) error {
	// Open or create the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Iterate over the map and write each key-value pair to the file
	for key, value := range m {
		_, err := fmt.Fprintf(file, "%v=%v\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
