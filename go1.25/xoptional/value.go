package xoptional

import "fmt"

func ValueWithOk[U any](opt T[U]) (U, bool) {
	if opt.HasValue() {
		return opt.Value(), true
	}
	var v U
	return v, false
}

// TODO: duplicate to xoptionalfb.OnEmpty
func ValueOr[U any](opt T[U], fb U) U {
	if opt.HasValue() {
		return opt.Value()
	}
	return fb
}

func ValueCastOr[U any, R any](opt T[U], cast func(U) R, fb R) R {
	if opt.HasValue() {
		return cast(opt.Value())
	}
	return fb
}

func ValueStringOr[U any](opt T[U], fb string) string {
	if opt.HasValue() {
		return fmt.Sprint(opt.Value())
	}
	return fb
}

func ValueOrSet[U any](opt T[U], setV U) U {
	if opt.HasValue() {
		return opt.Value()
	}
	opt.SetValue(setV)
	return setV
}

func ValueOrSetNew[U any](opt T[U], newV func() U) U {
	if opt.HasValue() {
		return opt.Value()
	}
	opt.SetValue(newV())
	return opt.Value()
}

func ValueOrTrySetNew[U any](opt T[U], newV func() (U, error)) (U, error) {
	if opt.HasValue() {
		return opt.Value(), nil
	}
	v, err := newV()
	if err != nil {
		return v, err
	}
	opt.SetValue(v)
	return v, nil
}
