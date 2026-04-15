package xslices

func Bind[T any, S ~[]T](s S) Slice[T] {
	return &slice[T]{s: s}
}

func BindComparable[T comparable, S ~[]T](s S) SliceOfComparable[T] {
	return &sliceCmp[T]{slice[T]{s: s}}
}

type Slice[T any] interface {
	// MapIn()
	// Map() Slice[T]
}

type SliceOfComparable[T comparable] interface {
	Slice[T]
	Has(v T) bool
}

type slice[T any] struct {
	s []T
}

type sliceCmp[T comparable] struct {
	slice[T]
}

func (s *sliceCmp[T]) Has(v T) bool {
	return Has(s.s, v)
}
