package core

import (
	"cmp"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"

	"github.com/Deimvis/go-ext/go1.25/xcheck/internal/buildtags"
)

// TODO: in debug mode format message even on success to validate that its passed properly
// TODO: EqAny (or AnyOf+pred?)

// NoErr checks whether err != nil and does nothing otherwise.
func NoErr(err error, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if err != nil {
		pcfg := printConfig{
			defaultBaseMsg: "err != nil",
			fmtValues: func() string {
				return err.Error()
			},
			showValues: true,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// True checks whether v != true and does nothing otherwise.
func True(v bool, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if !v {
		pcfg := printConfig{
			defaultBaseMsg: "not true",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NOTE: adding this may be a bad idea.
// Any boolean expression has a counterpart (like != for ==).
// So everything can be described with True() function
// and using True() everywhere makes code more uniform.
//
// False checks whether v != true and does nothing otherwise.
// func False(v bool, msgAndArgsAndOpts... any) (bool, string) {
// 	if v {
// 		msg := xfmt.Sprintfg(msgAndArgs...)
// 		return false, msg
// 	}
// 	return true, ""
// }

// NilPtr checks whether pointer v != nil and does nothing otherwise.
func NilPtr[T any](v *T, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v != nil {
		pcfg := printConfig{
			defaultBaseMsg: "not nil ptr",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NotNilPtr checks whether pointer v != nil and does nothing otherwise.
func NotNilPtr[T any](v *T, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v == nil {
		pcfg := printConfig{
			defaultBaseMsg: "nil ptr",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NilSlice checks whether slice v != nil and does nothing otherwise.
func NilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v != nil {
		pcfg := printConfig{
			defaultBaseMsg: "not nil slice",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NotNilSlice checks whether slice v != nil and does nothing otherwise.
func NotNilSlice[T any, S ~[]T](v S, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v == nil {
		pcfg := printConfig{
			defaultBaseMsg: "nil slice",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NilMap checks whether map v != nil and does nothing otherwise.
func NilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v != nil {
		pcfg := printConfig{
			defaultBaseMsg: "not nil map",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NotNilMap checks whether map v != nil and does nothing otherwise.
func NotNilMap[K comparable, V any, M ~map[K]V](v M, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v == nil {
		pcfg := printConfig{
			defaultBaseMsg: "nil map",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NilInterface checks whether interface v != nil and does nothing otherwise.
//
// Note that interface value being not nil does not guarantee that underlying value is not nil,
// example: `var x any = (*int)(nil); NilInterface(x) // false`
//
// Be careful since this function accepts any value,
// but check is valid only for interfaces (learn more in _cmp_archive.go).
func NilInterface(v any, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	buildtags.OnDebug(func() {
		rv := reflect.ValueOf(v)
		if v == nil && rv.Kind() != reflect.Invalid {
			fatal("xcheck: value is not interface, but %s\n%s", rv.Kind().String(), debug.Stack())
		}
	})

	if v != nil {
		pcfg := printConfig{
			defaultBaseMsg: "not nil interface",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NotNilInterface checks whether interface v != nil and does nothing otherwise.
//
// Note that interface value being not nil does not guarantee that underlying value is not nil,
// example: `var x any = (*int)(nil); NotNilInterface(x) // true`
//
// Be careful since this function accepts any value,
// but check is valid only for interfaces (learn more in _cmp_archive.go).
func NotNilInterface(v any, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	buildtags.OnDebug(func() {
		rv := reflect.ValueOf(v)
		if v == nil && rv.Kind() != reflect.Invalid {
			fatal("xcheck: value is not interface, but %s\n%s", rv.Kind().String(), debug.Stack())
		}
	})

	if v == nil {
		pcfg := printConfig{
			defaultBaseMsg: "nil interface",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// Nil checks whether v != nil and does nothing otherwise
//
// ISSUE: `any` value is nil iff it was directly assigned to nil.
// Casting any value to `any` interface would create a non-nil `any` variable.
// Example:
//
//	var myvar *int = nil
//	fmt.Println(myvar == nil) // true
//	var v any = myvar
//	fmt.Printn(v == nil) // false
func Nil(v any, msgAndArgsAndOpts ...any) (bool, string) {
	if v != nil {
		pcfg := printConfig{
			defaultBaseMsg: "not nil",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// NotNil checks whether v == nil and does nothing otherwise
func NotNil(v any, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v == nil {
		pcfg := printConfig{
			defaultBaseMsg: "is nil",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

func Eq[T comparable](v1 T, v2 T, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v1 != v2 {
		pcfg := printConfig{
			defaultBaseMsg: "not equal",
			fmtValues: func() string {
				// TODO: show diff when values are complex (slices, structs, etc)
				// TODO: format strings with quotes
				// TODO: add option to print type of values
				return fmt.Sprintf("%v != %v", v1, v2)
			},
			showValues: false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

func NotEq[T comparable](v1 T, v2 T, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if v1 == v2 {
		pcfg := printConfig{
			defaultBaseMsg: "equal",
			fmtValues: func() string {
				// TODO: show diff when values are complex (slices, structs, etc)
				return fmt.Sprintf("%v == %v", v1, v2)
			},
			showValues: false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

// TODO: forbid float in Lt
func Lt[T cmp.Ordered](v1 T, v2 T, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if !(v1 < v2) {
		pcfg := printConfig{
			defaultBaseMsg: "not less",
			fmtValues: func() string {
				// TODO: show diff when values are complex (slices, structs, etc)
				return fmt.Sprintf("%v >= %v", v1, v2)
			},
			showValues: false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

func Le[T cmp.Ordered](v1 T, v2 T, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	if !(v1 <= v2) {
		pcfg := printConfig{
			defaultBaseMsg: "not less or equal",
			fmtValues: func() string {
				// TODO: show diff when values are complex (slices, structs, etc)
				return fmt.Sprintf("%v > %v", v1, v2)
			},
			showValues: false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

func Implements[Target any](v any, msgAndArgsAndOpts ...any) (bool, string) {
	_validateMsg(msgAndArgsAndOpts...)
	_, ok := v.(Target)
	if !ok {
		pcfg := printConfig{
			defaultBaseMsg: "not implements",
			fmtValues: func() string {
				return fmt.Sprintf("%v not implements %T", v, *new(Target))
			},
			showValues: false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

var (
	fatal func(format string, v ...any) = log.Fatalf
)
