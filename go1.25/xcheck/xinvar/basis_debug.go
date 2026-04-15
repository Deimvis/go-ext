//go:build debug

package xinvar

import (
	"cmp"
	"log"

	"github.com/Deimvis/go-ext/go1.25/xcheck/internal/core"
)

// TODO: add print stack trace option to xcheck and enable by default for xinvar

// NoErr fails fast if err != nil and does nothing otherwise.
// Format string from msgAndArgs must contain last %s for error.
func NoErr(err error, msgAndArgs ...interface{}) {
	impl(core.NoErr(err, msgAndArgs...))
}

// True fails fast if v != true and does nothing otherwise.
func True(v bool, msgAndArgs ...interface{}) {
	impl(core.True(v, msgAndArgs...))
}

func NilPtr[T any](v *T, msgAndArgsAndOpts ...any) {
	impl(core.NilPtr(v, msgAndArgsAndOpts...))
}

func NotNilPtr[T any](v *T, msgAndArgsAndOpts ...any) {
	impl(core.NotNilPtr(v, msgAndArgsAndOpts...))
}

func NilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any) {
	impl(core.NilSlice(v, msgAndArgsAndOpts...))
}

func NotNilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any) {
	impl(core.NotNilSlice(v, msgAndArgsAndOpts...))
}

func NilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) {
	impl(core.NilMap(v, msgAndArgsAndOpts...))
}

func NotNilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) {
	impl(core.NotNilMap(v, msgAndArgsAndOpts...))
}

func NilInterface(v any, msgAndArgsAndOpts ...any) {
	impl(core.NilInterface(v, msgAndArgsAndOpts...))
}

func NotNilInterface(v any, msgAndArgsAndOpts ...any) {
	impl(core.NotNilInterface(v, msgAndArgsAndOpts...))
}

// Nil fails fast if v != nil and does nothing otherwise
func Nil(v any, msgAndArgs ...interface{}) {
	impl(core.Nil(v, msgAndArgs...))
}

// NotNil fails fast if v == nil and does nothing otherwise
func NotNil(v any, msgAndArgs ...interface{}) {
	impl(core.NotNil(v, msgAndArgs...))
}

func Eq[T comparable](v1 T, v2 T, msgAndArgs ...interface{}) {
	impl(core.Eq(v1, v2, msgAndArgs...))
}

func NotEq[T comparable](v1 T, v2 T, msgAndArgs ...interface{}) {
	impl(core.NotEq(v1, v2, msgAndArgs...))
}

func Lt[T cmp.Ordered](v1 T, v2 T, msgAndArgs ...interface{}) {
	impl(core.Lt(v1, v2, msgAndArgs...))
}

func Le[T cmp.Ordered](v1 T, v2 T, msgAndArgs ...interface{}) {
	impl(core.Le(v1, v2, msgAndArgs...))
}

func Implements[Target any](v any, msgAndArgs ...interface{}) {
	impl(core.Implements[Target](v))
}

func TypeImplements[T any, I any](msgAndArgs ...interface{}) {
	impl(core.TypeImplements[T, I](msgAndArgs...))
}

func impl(ok bool, errMsg string) {
	if !ok {
		exitFn(errMsg)
	}
}

var exitFn = func(msg string) {
	log.Fatalln(msg)
}
