package xchans

import "github.com/Deimvis/go-ext/go1.25/xslices"

// MapToSlice reads all elements from src and writes them to dst slice.
// It requires input channel to be closed in order to finish.
func MapToSlice[T any, U any, S ~[]U](src <-chan T, fn func(T) U, dst S) S {
	dst = dst[:0]
	for v := range src {
		dst = append(dst, fn(v))
	}
	return dst
}

// MapNToSlice reads N elements from src and writes them to dst slice.
// It requires input channel to pass through total of N elements or to be closed in order to finish.
// New dst slice will be allocated if cap(dst) < n.
// Number of items read equals len(dst).
func MapNToSlice[T any, U any, S ~[]U](src <-chan T, fn func(T) U, dst S, n int) S {
	dst = xslices.Resize(dst, n)
	for i := range n {
		v, ok := <-src
		if !ok {
			dst = dst[:i]
			break
		}
		dst[i] = fn(v)
	}
	return dst
}

func MapToNewSlice[T any, U any, S []U](src <-chan T, fn func(T) U) S {
	dst := make(S, newSliceDefaultCap)
	return MapToSlice(src, fn, dst)
}

func MapNToNewSlice[T any, U any, S []U](src <-chan T, fn func(T) U, n int) S {
	dst := make(S, newSliceDefaultCap)
	return MapNToSlice(src, fn, dst, n)
}

const (
	newSliceDefaultCap int = 0
)
