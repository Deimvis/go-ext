package xoptional

func FromPtr[U any](ptr *U) T[U] {
	o := T[U]{set: false}
	if ptr != nil {
		o.SetValue(*ptr)
	}
	return o
}

func ToPtr[U any](o T[U]) *U {
	var ptr *U = nil
	if o.HasValue() {
		v := o.Value()
		ptr = &v
	}
	return ptr
}

func ToPtrWithOk[U any](o T[U]) (*U, bool) {
	var ptr *U = nil
	var ok bool = false
	if o.HasValue() {
		v := o.Value()
		ptr = &v
		ok = true
	}
	return ptr, ok
}
