package xmaps

// Reverse reverses given map in a straightforward way.
// If map is not injective (doesn't form bijection),
// then output is non-deterministic.
func Reverse[K comparable, V comparable, M ~map[K]V, RevM map[V]K](m M) RevM {
	res := make(RevM, len(m))
	for k, v := range m {
		res[v] = k
	}
	return res
}
