package maps

import (
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test case 1: Map with string keys
	expected1 := []string{"a", "b", "c"}
	result1 := MapKeys(m)
	sort.Strings(result1) // Sort the keys
	assert.Equal(t, expected1, result1, "MapKeys(%v) = %v, expected %v", m, result1, expected1)

	// Test case 2: Map with int keys
	m2 := map[int]string{1: "one", 2: "two", 3: "three"}
	expected2 := []int{1, 2, 3}
	result2 := MapKeys(m2)
	sort.Ints(result2) // Sort the keys
	assert.Equal(t, expected2, result2, "MapKeys(%v) = %v, expected %v", m2, result2, expected2)

	// Test case 3: Empty map
	m3 := map[string]int{}
	expected3 := []string{}
	result3 := MapKeys(m3)
	sort.Strings(result3) // Sort the keys
	assert.Equal(t, expected3, result3, "MapKeys(%v) = %v, expected %v", m3, result3, expected3)
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test case 1: Map with int values
	expected1 := []int{1, 2, 3}
	result1 := MapValues(m)
	sort.Ints(result1) // Sort the values
	assert.Equal(t, expected1, result1, "MapValues(%v) = %v, expected %v", m, result1, expected1)

	// Test case 2: Map with string values
	m2 := map[int]string{1: "one", 2: "two", 3: "three"}
	expected2 := []string{"one", "two", "three"}
	result2 := MapValues(m2)
	sort.Strings(expected2)
	sort.Strings(result2) // Sort the values
	assert.Equal(t, expected2, result2, "MapValues(%v) = %v, expected %v", m2, result2, expected2)

	// Test case 3: Empty map
	m3 := map[string]int{}
	expected3 := []int{}
	result3 := MapValues(m3)
	sort.Ints(result3) // Sort the values
	assert.Equal(t, expected3, result3, "MapValues(%v) = %v, expected %v", m3, result3, expected3)
}
func TestMapContains(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test case 1: Key is present
	assert.True(t, MapContains(m, "a"), "MapContains(%v, %v) = false, expected true", m, "a")

	// Test case 2: Key is not present
	assert.False(t, MapContains(m, "d"), "MapContains(%v, %v) = true, expected false", m, "d")

	// Test case 3: Empty map
	m2 := map[string]int{}
	assert.False(t, MapContains(m2, "a"), "MapContains(%v, %v) = true, expected false", m2, "a")
}

func TestMapSortByKeys(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}

	// Test case 1: Sort by string keys in ascending order
	expected1 := map[string]int{"a": 1, "b": 2, "c": 3}
	result1 := MapSortByKeys(m, func(a, b string) bool { return a < b })
	assert.Equal(t, expected1, result1, "MapSortByKeys(%v) = %v, expected %v", m, result1, expected1)

	// Test case 2: Sort by string keys in descending order
	expected2 := map[string]int{"c": 3, "b": 2, "a": 1}
	result2 := MapSortByKeys(m, func(a, b string) bool { return a > b })
	assert.Equal(t, expected2, result2, "MapSortByKeys(%v) = %v, expected %v", m, result2, expected2)

	// Test case 3: Empty map
	m3 := map[string]int{}
	expected3 := map[string]int{}
	result3 := MapSortByKeys(m3, func(a, b string) bool { return a < b })
	assert.Equal(t, expected3, result3, "MapSortByKeys(%v) = %v, expected %v", m3, result3, expected3)
}

func TestMapSortByValues(t *testing.T) {
	m := map[string]int{"a": 3, "b": 1, "c": 2}

	// Test case 1: Sort by int values in ascending order
	expected1 := map[string]int{"b": 1, "c": 2, "a": 3}
	result1 := MapSortByValues(m, func(a, b int) bool { return a < b })
	assert.Equal(t, expected1, result1, "MapSortByValues(%v) = %v, expected %v", m, result1, expected1)

	// Test case 2: Sort by int values in descending order
	expected2 := map[string]int{"a": 3, "c": 2, "b": 1}
	result2 := MapSortByValues(m, func(a, b int) bool { return a > b })
	assert.Equal(t, expected2, result2, "MapSortByValues(%v) = %v, expected %v", m, result2, expected2)

	// Test case 3: Map with string values sorted in ascending order
	m2 := map[int]string{1: "three", 2: "one", 3: "two"}
	expected3 := map[int]string{2: "one", 3: "two", 1: "three"}
	result3 := MapSortByValues(m2, func(a, b string) bool { return a < b })
	assert.Equal(t, expected3, result3, "MapSortByValues(%v) = %v, expected %v", m2, result3, expected3)

	// Test case 4: Empty map
	m3 := map[string]int{}
	expected4 := map[string]int{}
	result4 := MapSortByValues(m3, func(a, b int) bool { return a < b })
	assert.Equal(t, expected4, result4, "MapSortByValues(%v) = %v, expected %v", m3, result4, expected4)
}

func TestMap2EnvFile(t *testing.T) {
	// Test case 1: Non-empty map
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	filePath := "/tmp/test_env_file_1.env"
	err := Map2EnvFile(m, filePath)
	assert.NoError(t, err, "Map2EnvFile(%v, %v) returned an error: %v", m, filePath, err)

	// Check if the file exists
	_, err = os.Stat(filePath)
	assert.False(t, os.IsNotExist(err), "File %v does not exist", filePath)

	// Test case 2: Empty map
	m2 := map[string]int{}
	filePath2 := "/tmp/test_env_file_2.env"
	err = Map2EnvFile(m2, filePath2)
	assert.NoError(t, err, "Map2EnvFile(%v, %v) returned an error: %v", m2, filePath2, err)

	// Check if the file exists
	_, err = os.Stat(filePath2)
	assert.False(t, os.IsNotExist(err), "File %v does not exist", filePath2)

	// Test case 3: Map with different types
	m3 := map[int]string{1: "one", 2: "two", 3: "three"}
	filePath3 := "/tmp/test_env_file_3.env"
	err = Map2EnvFile(m3, filePath3)
	assert.NoError(t, err, "Map2EnvFile(%v, %v) returned an error: %v", m3, filePath3, err)

	// Check if the file exists
	_, err = os.Stat(filePath3)
	assert.False(t, os.IsNotExist(err), "File %v does not exist", filePath3)

}
func TestMapFromSyncMap(t *testing.T) {
	// Test case 1: SyncMap with string keys and int values
	sm1 := &sync.Map{}
	sm1.Store("a", 1)
	sm1.Store("b", 2)
	sm1.Store("c", 3)

	expected1 := map[string]int{"a": 1, "b": 2, "c": 3}
	result1 := MapFromSyncMap[string, int](sm1)
	assert.Equal(t, expected1, result1, "MapFromSyncMap(%v) = %v, expected %v", sm1, result1, expected1)

	// Test case 2: SyncMap with int keys and string values
	sm2 := &sync.Map{}
	sm2.Store(1, "one")
	sm2.Store(2, "two")
	sm2.Store(3, "three")

	expected2 := map[int]string{1: "one", 2: "two", 3: "three"}
	result2 := MapFromSyncMap[int, string](sm2)
	assert.Equal(t, expected2, result2, "MapFromSyncMap(%v) = %v, expected %v", sm2, result2, expected2)

	// Test case 3: Empty SyncMap
	sm3 := &sync.Map{}
	expected3 := map[string]int{}
	result3 := MapFromSyncMap[string, int](sm3)
	assert.Equal(t, expected3, result3, "MapFromSyncMap(%v) = %v, expected %v", sm3, result3, expected3)
}
