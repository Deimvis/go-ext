package xcheck

import "github.com/Deimvis/go-ext/go1.25/xcheck/internal/core"

func Eq[T comparable]() core.EqPred[T] {
	return core.EqPred[T]{}
}
