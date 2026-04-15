package xreflect

import (
	"reflect"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xinvar"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

// FillNilPointers recursively fills nil pointers with pointers to
// default value of the underlying type.
// Essentially, it makes structure pointers dereference safe.
// Input value must be a non-nil pointer to a struct or a non-nil pointer to a pointer to a struct,
// in the latter case this struct will be filled as well.
func FillNilPointers(s any) {
	opts := fillNilPointersOptions{
		onlyStructPtrs: false,
	}
	fillNilPointers(s, opts)
}

// FillNilStructPointers works same as FillNilPointers,
// but fills only nil pointers to structs (other pointers will remain nil).
func FillNilStructPointers(s any) {
	opts := fillNilPointersOptions{
		onlyStructPtrs: true,
	}
	fillNilPointers(s, opts)
}

func fillNilPointers(s any, opts fillNilPointersOptions) {
	v := reflect.ValueOf(s)
	xmust.Eq(v.Kind(), reflect.Pointer)
	xmust.True(!v.IsNil())
	if v.Elem().Kind() == reflect.Pointer {
		structType := v.Elem().Type().Elem()
		xmust.Eq(structType.Kind(), reflect.Struct, "got **<non-struct> (pointer to a pointer to not struct)")
		if v.Elem().IsNil() {
			newStruct := reflect.New(structType)
			v.Elem().Set(newStruct)
		}
		v = v.Elem()
		xinvar.Eq(v.Elem().Kind(), reflect.Struct)
	}
	xmust.Eq(v.Elem().Kind(), reflect.Struct, "got neither *<struct> nor **<struct>")
	fillPointerFields(v.Elem(), opts)
}

func fillPointerFields(v reflect.Value, opts fillNilPointersOptions) {
	xinvar.Eq(v.Kind(), reflect.Struct)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldV := v.Field(i)
		fieldInfo := t.Field(i)

		if !fieldV.CanSet() {
			continue
		}

		switch fieldInfo.Type.Kind() {
		case reflect.Pointer:
			if opts.onlyStructPtrs && fieldInfo.Type.Elem().Kind() != reflect.Struct {
				break
			}
			if fieldV.IsNil() {
				newValue := reflect.New(fieldInfo.Type.Elem())
				fieldV.Set(newValue)
			}
			if fieldInfo.Type.Elem().Kind() != reflect.Struct {
				break
			}

			fieldV = fieldV.Elem()
			fallthrough

		case reflect.Struct:
			fillPointerFields(fieldV, opts)

		case reflect.Slice:
			if opts.onlyStructPtrs {
				break
			}
			if fieldV.IsNil() {
				fieldV.Set(reflect.MakeSlice(fieldInfo.Type, 0, 0))
			}

		case reflect.Map:
			if opts.onlyStructPtrs {
				break
			}
			if fieldV.IsNil() {
				fieldV.Set(reflect.MakeMap(fieldInfo.Type))
			}
		}
	}
}

type fillNilPointersOptions struct {
	onlyStructPtrs bool
}

// resolve recursively resolves pointers and interfaces to their underlying value.
// https://github.com/go-playground/validator/blob/a947377040f8ebaee09f20d09a745ec369396793/util.go#L15
func resolve(v reflect.Value) reflect.Value {

BEGIN:
	switch v.Kind() {
	case reflect.Pointer:

		if v.IsNil() {
			return v
		}

		v = v.Elem()
		goto BEGIN

	case reflect.Interface:

		if v.IsNil() {
			return v
		}

		v = v.Elem()
		goto BEGIN

	case reflect.Invalid:
		return v

	default:
		return v
	}
}
