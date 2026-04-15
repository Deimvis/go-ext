package xslices

import (
	"maps"

	"github.com/Deimvis/go-ext/go1.25/xmaps"
)

// TODO: add option for unique count estimation (to preallocate map
// and so work more efficiently with slices containing a lot of unique items)
func ToInclusionMap[T comparable, S ~[]T](s S) xmaps.InclusionMap[T] {
	it := func(yield func(k T, v struct{}) bool) {
		for _, v := range s {
			if !yield(v, struct{}{}) {
				return
			}
		}
	}
	return maps.Collect(it)
}
