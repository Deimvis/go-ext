package xoptional

import (
	"errors"
	"fmt"
)

func New[U any](v ...U) T[U] {
	if len(v) > 1 {
		panic(errors.New("multiple values"))
	}
	o := T[U]{set: false}
	if len(v) == 1 {
		o.SetValue(v[0])
	}
	return o
}

// T (optional) is imlemented
// not using `type optional.T[U any] *U`
// because this way Golang doesn't allow
// to define methods.
// Also, we want to follow user's intent -
// if one uses value U, then he expects
// to move around U and not *U.
// If type U is large and you would like to move it
// by pointer, then use optional.T[*U] and not *optional.T[U].
type T[U any] struct {
	v   U
	set bool
}

func (o T[U]) HasValue() bool {
	return o.set
}

func (o T[U]) Value() U {
	// in order to prevent accident mistakes
	if !o.set {
		panic(errNoValue)
	}
	return o.v
}

// ValuePtr is guaranteed to return
// nil pointer when HasValue() is false.
func (o *T[U]) ValuePtr() *U {
	// in order to prevent accident mistakes
	// when value was not set exlicitly
	// but was modified by pointer
	// and latter code works unexpectedly
	// since HasValue() is still false.
	if !o.set {
		return nil
	}
	return &o.v
}

func (o *T[U]) SetValue(v U) {
	o.v = v
	o.set = true
}

func (o *T[U]) Reset() {
	o.set = false
}

// String is useful for debugging,
// do not rely on its value.
// xoptional.ValueStringOr is preferred
// for more consistent formatting behaviour.
func (o T[U]) String() string {
	if o.set {
		return fmt.Sprint(o.v)
	}
	return "nil"
}

var errNoValue = errors.New("no value")

// TODO: impl
// type AtomicT[U any] struct {
// 	v atomic.Pointer[struct {
// 		v   U
// 		set bool
// 	}]
// }
