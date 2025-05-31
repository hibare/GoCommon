// Package slice provides utilities for working with slices.
package slice

// Unique returns a slice with unique elements.
func Unique[T comparable](list []T) []T {
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

// Diff returns the difference between two slices.
func Diff[T comparable](a, b []T) []T {
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

// Intersect returns the intersection between two slices.
func Intersect[T comparable](a, b []T) []T {
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

// Union returns the union between two slices.
func Union[T comparable](a, b []T) []T {
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
