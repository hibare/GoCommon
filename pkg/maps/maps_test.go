package maps

import (
	"sort"
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
