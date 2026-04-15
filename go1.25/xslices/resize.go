package xslices

// Resize resizes the slice to the specified length,
// possibly modifying original slice.
func Resize[T any, S ~[]T](s S, length int) S {
	if length <= cap(s) {
		return s[:length]
	}

	snew := make(S, length)
	copy(snew, s)
	return snew
}
