package xmapfb

func NoKey[K comparable, V any, M ~map[K]V](m M, k K, fb V) V {
	if v, ok := m[k]; ok {
		return v
	}
	return fb
}
