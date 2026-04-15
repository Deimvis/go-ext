package xslices

import (
	"cmp"
	"sort"

	"github.com/Deimvis/go-ext/go1.25/xfn"
)

// TODO: research inplace function signatures. accepting pointer makes code more readable since you pass a pointer
// and therefore hint readers that object will be modified, but it prohibits from chaining functions, like
// SortIn(DeduplicateIn(s)).
// Maybe there should be another abstraction for chaining like:
// xslices.Pipeline(s, xslices.SortOp, xslices.DeduplicateOp)
// or
// xslicepl.Run(s, xslicepl.SortOp, xslicepl.DeduplicateOp)

// TODO: sort options: Sort(s, WithKeyFn(...)), Sort(s, WithCmpFn(func(a, b T) int))
//       maybe move options to xslicesort: xslices.Sort(s, xlicesort.ByKey(MyType.GetTs), xslicesort.Reversed())
// TODO: stable sort

func Sort[T cmp.Ordered, S ~[]T](s S) S {
	return SortFn(s, xfn.Id)
}

func SortIn[T cmp.Ordered, S ~[]T](s *S) S {
	return SortFnIn(s, xfn.Id)
}

func SortFn[T any, U cmp.Ordered, S ~[]T](s S, keyFn func(v T) U) S {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return SortFnIn(&scopy, keyFn)
}

func SortFnIn[T any, U cmp.Ordered, S ~[]T](s *S, keyFn func(v T) U) S {
	sort.SliceStable(*s, func(i int, j int) bool {
		return keyFn((*s)[i]) < keyFn((*s)[j])
	})
	return *s
}
