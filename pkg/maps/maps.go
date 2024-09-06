package maps

import "sort"

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
