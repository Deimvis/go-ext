package xiter

import "iter"

type ReduceFn[T any] func(T, T) T
type Reduce2Fn[T1, T2 any] func(T1, T2, T1, T2) (T1, T2)

func Reduce[T any](seq iter.Seq[T], fn ReduceFn[T], init T) T {
	res := init
	seq(func(v T) bool {
		res = fn(res, v)
		return true
	})
	return res
}
