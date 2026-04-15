package xslices

import (
	"iter"
)

// StopBefore stops when slice reaches item equal stopV
// and does NOT return stopV.
// If you need to stop at particular index, just
// make a new slice: s = s[:stopInd] =)
func StopBefore[T comparable, S ~[]T](s S, stopV T) iter.Seq2[int, T] {
	return func(yield func(i int, v T) bool) {
		for i, v := range s {
			if v == stopV {
				return
			}
			if !yield(i, v) {
				return
			}
		}
	}
}
