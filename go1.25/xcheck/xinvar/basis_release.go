//go:build !debug

package xinvar

import "cmp"

func NoErr(err error, msgAndArgs ...interface{})                               {}
func True(v bool, msgAndArgs ...interface{})                                   {}
func NilPtr[T any](v *T, msgAndArgsAndOpts ...any)                             {}
func NotNilPtr[T any](v *T, msgAndArgsAndOpts ...any)                          {}
func NilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any)                    {}
func NotNilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any)                 {}
func NilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any)    {}
func NotNilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) {}
func NilInterface(v any, msgAndArgsAndOpts ...any)                             {}
func NotNilInterface(v any, msgAndArgsAndOpts ...any)                          {}
func Nil(v any, msgAndArgs ...interface{})                                     {}
func NotNil(v any, msgAndArgs ...interface{})                                  {}
func Eq[T comparable](v1 T, v2 T, msgAndArgs ...interface{})                   {}
func NotEq[T comparable](v1 T, v2 T, msgAndArgs ...interface{})                {}
func Lt[T cmp.Ordered](v1 T, v2 T, msgAndArgs ...interface{})                  {}
func Le[T cmp.Ordered](v1 T, v2 T, msgAndArgs ...interface{})                  {}
func Implements[Target any](v any, msgAndArgs ...interface{})                  {}
func TypeImplements[T any, I any](msgAndArgs ...interface{})                   {}
