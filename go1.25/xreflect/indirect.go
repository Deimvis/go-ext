package xreflect

import "reflect"

// Indirect takes Elem() of value
// if value's kind allows it
// and value is not nil.
// Safe to call with value of any kind.
// Returns true if indirection happened.
func Indirect(v reflect.Value) (reflect.Value, bool) {
	switch v.Kind() {
	case reflect.Pointer:
		return IndirectPointer(v)
	case reflect.Interface:
		return IndirectInterface(v)
	}
	return v, false
}

// IndirectPointer takes Elem()
// of value with kind Pointer
// if value is not nil.
// Behaviour is undefined if input
// value's kind is not Poitner.
// Returns true if indirection happened.
func IndirectPointer(v reflect.Value) (reflect.Value, bool) {
	if v.IsNil() {
		return v, false
	}
	return v.Elem(), true
}

// IndirectInterface takes Elem()
// of value with kind Interface
// if value is not nil.
// Behaviour is undefined if input
// value's kind is not Interface.
// Returns true if indirection happened.
func IndirectInterface(v reflect.Value) (reflect.Value, bool) {
	if v.IsNil() {
		return v, false
	}
	return v.Elem(), true
}

// RecursiveIndirect recursively calls Indirect
// on value until Indirect returns false.
func RecursiveIndirect(v reflect.Value) reflect.Value {
	ok := true
	for ok {
		v, ok = Indirect(v)
	}
	return v
}

// TODO: implement ReportRecursiveIndirect,
// which returns report on how recursive indirect happened
// (describes the whole indirection process which happened)
// func ReportRecursiveIndirect(v reflect.Value) (reflect.Value, Report) {
// }
