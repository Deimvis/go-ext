package xiter

import "iter"

func Map[T any, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// func Map2[T1, T2 any, U1, U2 any](seq iter.Seq2[T1, T2], fn func(T1, T2) (U1, U2)) iter.Seq2[U1, U2] {
// }

// func Map12[T any, U1, U2 any](seq iter.Seq[T], fn func(T) (U1, U2)) iter.Seq2[U1, U2] {}
// func Map21[T1, T2 any, U any](seq iter.Seq2[T1, T2], fn func(T1, T2) U) iter.Seq[U] {}
