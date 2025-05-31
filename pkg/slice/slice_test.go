package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceUnique(t *testing.T) {
	require.ElementsMatch(t, []string{"1", "2", "3"}, Unique([]string{"1", "2", "3", "1", "2"}))
	require.ElementsMatch(t, []int{1, 2, 3}, Unique([]int{1, 2, 3, 1, 2}))
	require.ElementsMatch(t, []float64{1.1, 2.2, 3.3}, Unique([]float64{1.1, 2.2, 3.3, 1.1, 2.2}))
	require.ElementsMatch(t, []string{}, Unique([]string{}))
	require.ElementsMatch(t, []int{1}, Unique([]int{1, 1, 1, 1}))
}

func TestSliceDiff(t *testing.T) {
	require.ElementsMatch(t, []string{"1"}, Diff([]string{"1", "2", "3"}, []string{"2", "3"}))
	require.ElementsMatch(t, []int{1}, Diff([]int{1, 2, 3}, []int{2, 3}))
	require.ElementsMatch(t, []float64{1.1}, Diff([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3}))
	require.ElementsMatch(t, []string{"1", "2", "3"}, Diff([]string{"1", "2", "3"}, []string{}))
	require.ElementsMatch(t, []int{}, Diff([]int{1, 2, 3}, []int{1, 2, 3}))
	require.ElementsMatch(t, []int{}, Diff([]int{}, []int{1, 2, 3}))
}

func TestSliceIntersect(t *testing.T) {
	require.ElementsMatch(t, []string{"2", "3"}, Intersect([]string{"1", "2", "3"}, []string{"2", "3", "4"}))
	require.ElementsMatch(t, []int{2, 3}, Intersect([]int{1, 2, 3}, []int{2, 3, 4}))
	require.ElementsMatch(t, []float64{2.2, 3.3}, Intersect([]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3, 4.4}))
	require.ElementsMatch(t, []string{}, Intersect([]string{"1", "2", "3"}, []string{"4", "5", "6"}))
	require.ElementsMatch(t, []int{}, Intersect([]int{1, 2, 3}, []int{}))
	require.ElementsMatch(t, []int{}, Intersect([]int{}, []int{1, 2, 3}))
}

func TestSliceUnion(t *testing.T) {
	require.ElementsMatch(t, []string{"1", "2", "3", "4"}, Union([]string{"1", "2", "3"}, []string{"3", "4"}))
	require.ElementsMatch(t, []int{1, 2, 3, 4}, Union([]int{1, 2, 3}, []int{3, 4}))
	require.ElementsMatch(t, []float64{1.1, 2.2, 3.3, 4.4}, Union([]float64{1.1, 2.2, 3.3}, []float64{3.3, 4.4}))
	require.ElementsMatch(t, []string{"1", "2", "3"}, Union([]string{"1", "2", "3"}, []string{}))
	require.ElementsMatch(t, []int{1, 2, 3}, Union([]int{1, 2, 3}, []int{}))
	require.ElementsMatch(t, []int{1, 2, 3}, Union([]int{}, []int{1, 2, 3}))
	require.ElementsMatch(t, []string{}, Union([]string{}, []string{}))
}
