package xslices

// Fill fills in-place given slice
// with values constructed using newFn.
func Fill[T any, S ~[]T](s *S, newFn func(int) T) {
	for i := range len(*s) {
		(*s)[i] = newFn(i)
	}
}

// FillConst fills in-place given slice
// with given const value.
func FillConst[T any, S ~[]T](s *S, v T) {
	for i := range len(*s) {
		(*s)[i] = v
	}
}
