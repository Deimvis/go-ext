package xslices

// GroupBy groups elements of given slice by key derived using given keyFn.
// Order of elements in each bucket is not guaranteed.
func GroupBy[T any, K comparable](s []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T, len(s))
	for _, x := range s {
		key := keyFn(x)
		result[key] = append(result[key], x)
	}
	return result
}

// GroupIndBy groups element indices of given slice by key derived using given keyFn.
// Order of elements in each bucket is not guaranteed.
func GroupIndBy[T any, K comparable](s []T, keyFn func(T) K) map[K][]int {
	result := make(map[K][]int, len(s))
	for i, x := range s {
		key := keyFn(x)
		result[key] = append(result[key], i)
	}
	return result
}
