package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, SliceContains("1", []string{"1", "2", "3"}))
	assert.False(t, SliceContains("11", []string{"1", "2", "3"}))
	assert.True(t, SliceContains(2, []int{1, 2, 3}))
	assert.False(t, SliceContains(22, []int{1, 2, 3}))
	assert.True(t, SliceContains(3.3, []float64{1.1, 2.2, 3.3}))
	assert.False(t, SliceContains(33.3, []float64{1.1, 2.2, 3.3}))
}

func TestSliceUnique(t *testing.T) {
	assert.ElementsMatch(t, []string{"1", "2", "3"}, SliceUnique([]string{"1", "2", "3", "1", "2"}))
	assert.ElementsMatch(t, []int{1, 2, 3}, SliceUnique([]int{1, 2, 3, 1, 2}))
	assert.ElementsMatch(t, []float64{1.1, 2.2, 3.3}, SliceUnique([]float64{1.1, 2.2, 3.3, 1.1, 2.2}))
	assert.ElementsMatch(t, []string{}, SliceUnique([]string{}))
	assert.ElementsMatch(t, []int{1}, SliceUnique([]int{1, 1, 1, 1}))
}

func TestSliceDiff(t *testing.T) {
	assert.ElementsMatch(t, []string{"1"}, SliceDiff([]string{"1", "2", "3"}, []string{"2", "3"}))
	assert.ElementsMatch(t, []int{1}, SliceDiff([]int{1, 2, 3}, []int{2, 3}))
	assert.ElementsMatch(t, []float64{1.1}, SliceDiff([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3}))
	assert.ElementsMatch(t, []string{"1", "2", "3"}, SliceDiff([]string{"1", "2", "3"}, []string{}))
	assert.ElementsMatch(t, []int{}, SliceDiff([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.ElementsMatch(t, []int{}, SliceDiff([]int{}, []int{1, 2, 3}))
}

func TestSliceIntersect(t *testing.T) {
	assert.ElementsMatch(t, []string{"2", "3"}, SliceIntersect([]string{"1", "2", "3"}, []string{"2", "3", "4"}))
	assert.ElementsMatch(t, []int{2, 3}, SliceIntersect([]int{1, 2, 3}, []int{2, 3, 4}))
	assert.ElementsMatch(t, []float64{2.2, 3.3}, SliceIntersect([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3, 4.4}))
	assert.ElementsMatch(t, []string{}, SliceIntersect([]string{"1", "2", "3"}, []string{"4", "5", "6"}))
	assert.ElementsMatch(t, []int{}, SliceIntersect([]int{1, 2, 3}, []int{}))
	assert.ElementsMatch(t, []int{}, SliceIntersect([]int{}, []int{1, 2, 3}))
}

func TestSliceUnion(t *testing.T) {
	assert.ElementsMatch(t, []string{"1", "2", "3", "4"}, SliceUnion([]string{"1", "2", "3"}, []string{"3", "4"}))
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, SliceUnion([]int{1, 2, 3}, []int{3, 4}))
	assert.ElementsMatch(t, []float64{1.1, 2.2, 3.3, 4.4}, SliceUnion([]float64{1.1, 2.2, 3.3}, []float64{3.3, 4.4}))
	assert.ElementsMatch(t, []string{"1", "2", "3"}, SliceUnion([]string{"1", "2", "3"}, []string{}))
	assert.ElementsMatch(t, []int{1, 2, 3}, SliceUnion([]int{1, 2, 3}, []int{}))
	assert.ElementsMatch(t, []int{1, 2, 3}, SliceUnion([]int{}, []int{1, 2, 3}))
	assert.ElementsMatch(t, []string{}, SliceUnion([]string{}, []string{}))
}
