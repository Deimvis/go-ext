package xmust

import (
	"errors"

	"github.com/Deimvis/go-ext/go1.25/xfmt"
)

func MustEqualInt(v1 int, v2 int, msgAndArgs ...interface{}) {
	EqualComparable(v1, v2, msgAndArgs...)
}

func EqualComparable[T comparable](v1 T, v2 T, msgAndArgs ...interface{}) {
	if v1 != v2 {
		msg := xfmt.Sprintfg(msgAndArgs...)
		panic(errors.New(msg))
	}
}
