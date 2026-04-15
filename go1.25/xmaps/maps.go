package xmaps

// TODO: get rid of obsolete functions for go1.21 and go1.23

// CloneMap returns a shallow copy of an original map
func Clone[K comparable, V any](m map[K]V) map[K]V {
	mcopy := make(map[K]V)
	for k, v := range m {
		mcopy[k] = v
	}
	return mcopy
}

// Map maps a map :)
func Map[K comparable, V any, U any](m map[K]V, mapFn func(k K, v V) (K, U)) map[K]U {
	res := make(map[K]U, len(m))
	for k, v := range m {
		k1, v1 := mapFn(k, v)
		res[k1] = v1
	}
	return res
}

func Filter[K comparable, V any](m map[K]V, filFn func(k K, v V) bool) map[K]V {
	mcopy := Clone(m)
	FilterIn(&mcopy, filFn)
	return mcopy
}

// FilteIn filters map in-place.
// Keeps only those elements for which given filFn returns true.
func FilterIn[K comparable, V any](m *map[K]V, filFn func(k K, v V) bool) {
	for k, v := range *m {
		if !filFn(k, v) {
			delete(*m, k)
		}
	}
}

func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, len(m))
	i := 0
	for k := range m {
		res[i] = k
		i += 1
	}
	return res
}

func Values[K comparable, V any](m map[K]V) []V {
	res := make([]V, len(m))
	i := 0
	for _, v := range m {
		res[i] = v
		i += 1
	}
	return res
}

func HasKey[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

// MergeIn merges right map into left, modifying left map.
func MergeIn[K comparable, V any](left map[K]V, right map[K]V) map[K]V {
	for k, v := range right {
		left[k] = v
	}
	return left
}

func Merge[K comparable, V any](left map[K]V, right map[K]V) map[K]V {
	res := make(map[K]V)
	for k, v := range left {
		res[k] = v
	}
	for k, v := range right {
		res[k] = v
	}
	return res
}
