// Package slice provides utilities for working with slices.
package slice

// StringInSlice checks if a string is present in slice.
func SliceContains[T comparable](a T, list []T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SliceUnique returns a slice with unique elements.
func SliceUnique[T comparable](list []T) []T {
	keys := make(map[T]bool)
	listUnique := []T{}
	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			listUnique = append(listUnique, entry)
		}
	}
	return listUnique
}

// SliceDiff returns the difference between two slices.
func SliceDiff[T comparable](a, b []T) []T {
	mb := make(map[T]bool)
	for _, x := range b {
		mb[x] = true
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// SliceIntersect returns the intersection between two slices.
func SliceIntersect[T comparable](a, b []T) []T {
	mb := make(map[T]bool)
	for _, x := range b {
		mb[x] = true
	}
	var intersect []T
	for _, x := range a {
		if _, found := mb[x]; found {
			intersect = append(intersect, x)
		}
	}
	return intersect
}

// SliceUnion returns the union between two slices.
func SliceUnion[T comparable](a, b []T) []T {
	ma := make(map[T]bool)
	for _, x := range a {
		ma[x] = true
	}
	for _, x := range b {
		ma[x] = true
	}
	var union []T
	for x := range ma {
		union = append(union, x)
	}
	return union
}

// SliceIndexOf returns the index of the first occurrence of a in list, or -1 if not found.
func SliceIndexOf[T comparable](a T, list []T) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}
