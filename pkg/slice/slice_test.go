package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, Contains("1", []string{"1", "2", "3"}))
	assert.False(t, Contains("11", []string{"1", "2", "3"}))
	assert.True(t, Contains(2, []int{1, 2, 3}))
	assert.False(t, Contains(22, []int{1, 2, 3}))
	assert.True(t, Contains(3.3, []float64{1.1, 2.2, 3.3}))
	assert.False(t, Contains(33.3, []float64{1.1, 2.2, 3.3}))
}

func TestSliceUnique(t *testing.T) {
	assert.ElementsMatch(t, []string{"1", "2", "3"}, Unique([]string{"1", "2", "3", "1", "2"}))
	assert.ElementsMatch(t, []int{1, 2, 3}, Unique([]int{1, 2, 3, 1, 2}))
	assert.ElementsMatch(t, []float64{1.1, 2.2, 3.3}, Unique([]float64{1.1, 2.2, 3.3, 1.1, 2.2}))
	assert.ElementsMatch(t, []string{}, Unique([]string{}))
	assert.ElementsMatch(t, []int{1}, Unique([]int{1, 1, 1, 1}))
}

func TestSliceDiff(t *testing.T) {
	assert.ElementsMatch(t, []string{"1"}, Diff([]string{"1", "2", "3"}, []string{"2", "3"}))
	assert.ElementsMatch(t, []int{1}, Diff([]int{1, 2, 3}, []int{2, 3}))
	assert.ElementsMatch(t, []float64{1.1}, Diff([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3}))
	assert.ElementsMatch(t, []string{"1", "2", "3"}, Diff([]string{"1", "2", "3"}, []string{}))
	assert.ElementsMatch(t, []int{}, Diff([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.ElementsMatch(t, []int{}, Diff([]int{}, []int{1, 2, 3}))
}

func TestSliceIntersect(t *testing.T) {
	assert.ElementsMatch(t, []string{"2", "3"}, Intersect([]string{"1", "2", "3"}, []string{"2", "3", "4"}))
	assert.ElementsMatch(t, []int{2, 3}, Intersect([]int{1, 2, 3}, []int{2, 3, 4}))
	assert.ElementsMatch(t, []float64{2.2, 3.3}, Intersect([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3, 4.4}))
	assert.ElementsMatch(t, []string{}, Intersect([]string{"1", "2", "3"}, []string{"4", "5", "6"}))
	assert.ElementsMatch(t, []int{}, Intersect([]int{1, 2, 3}, []int{}))
	assert.ElementsMatch(t, []int{}, Intersect([]int{}, []int{1, 2, 3}))
}

func TestSliceUnion(t *testing.T) {
	assert.ElementsMatch(t, []string{"1", "2", "3", "4"}, Union([]string{"1", "2", "3"}, []string{"3", "4"}))
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, Union([]int{1, 2, 3}, []int{3, 4}))
	assert.ElementsMatch(t, []float64{1.1, 2.2, 3.3, 4.4}, Union([]float64{1.1, 2.2, 3.3}, []float64{3.3, 4.4}))
	assert.ElementsMatch(t, []string{"1", "2", "3"}, Union([]string{"1", "2", "3"}, []string{}))
	assert.ElementsMatch(t, []int{1, 2, 3}, Union([]int{1, 2, 3}, []int{}))
	assert.ElementsMatch(t, []int{1, 2, 3}, Union([]int{}, []int{1, 2, 3}))
	assert.ElementsMatch(t, []string{}, Union([]string{}, []string{}))
}

func TestSliceIndexOf(t *testing.T) {
	assert.Equal(t, 0, IndexOf("1", []string{"1", "2", "3"}))
	assert.Equal(t, -1, IndexOf("11", []string{"1", "2", "3"}))
	assert.Equal(t, 1, IndexOf(2, []int{1, 2, 3}))
	assert.Equal(t, -1, IndexOf(22, []int{1, 2, 3}))
	assert.Equal(t, 2, IndexOf(3.3, []float64{1.1, 2.2, 3.3}))
	assert.Equal(t, -1, IndexOf(33.3, []float64{1.1, 2.2, 3.3}))
	assert.Equal(t, -1, IndexOf("1", []string{}))
}
