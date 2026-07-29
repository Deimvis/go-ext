package xslices

// TODO: fix Copy(nil) != nil
// TODO: remove since slices.Clone now exists
// - (it is as efficient as this)
func Copy[T any, S ~[]T](s S) S {
	sCopy := make([]T, len(s))
	copy(sCopy, s)
	return sCopy
}
