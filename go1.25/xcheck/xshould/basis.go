package xshould

import (
	"cmp"
	"errors"

	"github.com/Deimvis/go-ext/go1.25/xcheck/internal/core"
)

func NoErr(err error, msgAndArgs ...interface{}) error {
	return impl(core.NoErr(err, msgAndArgs...))
}

func True(v bool, msgAndArgs ...interface{}) error {
	return impl(core.True(v, msgAndArgs...))
}

func NilPtr[T any](v *T, msgAndArgsAndOpts ...any) error {
	return impl(core.NilPtr(v, msgAndArgsAndOpts...))
}

func NotNilPtr[T any](v *T, msgAndArgsAndOpts ...any) error {
	return impl(core.NotNilPtr(v, msgAndArgsAndOpts...))
}

func NilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any) error {
	return impl(core.NilSlice(v, msgAndArgsAndOpts...))
}

func NotNilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any) error {
	return impl(core.NotNilSlice(v, msgAndArgsAndOpts...))
}

func NilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) error {
	return impl(core.NilMap(v, msgAndArgsAndOpts...))
}

func NotNilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) error {
	return impl(core.NotNilMap(v, msgAndArgsAndOpts...))
}

func NilInterface(v any, msgAndArgsAndOpts ...any) error {
	return impl(core.NilInterface(v, msgAndArgsAndOpts...))
}

func NotNilInterface(v any, msgAndArgsAndOpts ...any) error {
	return impl(core.NotNilInterface(v, msgAndArgsAndOpts...))
}

func Nil(v any, msgAndArgs ...interface{}) error {
	return impl(core.Nil(v, msgAndArgs...))
}

func NotNil(v any, msgAndArgs ...interface{}) error {
	return impl(core.NotNil(v, msgAndArgs...))
}

func Eq[T comparable](v1 T, v2 T, msgAndArgs ...interface{}) error {
	return impl(core.Eq(v1, v2, msgAndArgs...))
}

func NotEq[T comparable](v1 T, v2 T, msgAndArgs ...interface{}) error {
	return impl(core.NotEq(v1, v2, msgAndArgs...))
}

func AnyOf[T any](pred core.BinaryPredicate_[T], v1 T, v2opts []T, msgAndArgs ...interface{}) error {
	return impl(core.AnyOf(pred, v1, v2opts, msgAndArgs...))
}

func Lt[T cmp.Ordered](v1 T, v2 T, msgAndArgs ...interface{}) error {
	return impl(core.Lt(v1, v2, msgAndArgs...))
}

func Le[T cmp.Ordered](v1 T, v2 T, msgAndArgs ...interface{}) error {
	return impl(core.Le(v1, v2, msgAndArgs...))
}

func Implements[Target any](v any, msgAndArgs ...interface{}) error {
	return impl(core.Implements[Target](v, msgAndArgs...))
}

func TypeImplements[T any, I any](msgAndArgs ...interface{}) error {
	return impl(core.TypeImplements[T, I](msgAndArgs...))
}

func impl(ok bool, errMsg string) error {
	if !ok {
		return errors.New(errMsg)
	}
	return nil
}
