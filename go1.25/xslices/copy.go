package xslices

func Copy[T any, S ~[]T](s S) S {
	sCopy := make([]T, len(s))
	copy(sCopy, s)
	return sCopy
}
