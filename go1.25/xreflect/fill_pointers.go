package xreflect

import (
	"reflect"

	"github.com/Deimvis/go-ext/go1.25/xcheck"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xinvar"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

// TODO: consider renaming to FillNilFieldPointers, or somehow canonizing.
// * opt) FillFields(v, xreflect.NilPtr)
//        FillFields(v, xreflect.NilStructPtr)
// * opt) FillNilPtrFields(v) + FillNilStructPtrFields(v)
// * opt) FillNilPtrFields(v, xreflectfill.OnlyStructPtrs)
// * opt) FillNilPointers, but make work with slices and other containers, but make work with filters (only struct pointers)
// * opt) FillNils(v, Pointers, Slices, Maps) // works with input container or struct, checks CanAddr; by default fills any nils (except interfaces)

// FillNilPointers recursively fills nil pointers with pointers to
// default value of the underlying type.
// Essentially, it makes structure pointers dereference safe.
// Input value must be a non-nil pointer to a struct or a non-nil pointer to a pointer to a struct,
// in the latter case this struct will be filled as well.
func FillNilPointers(s any) {
	opts := fillNilPointersOptions{
		onlyStructPtrs: false,
	}
	fillNilPointers(reflect.ValueOf(s), opts)
}

// FillNilStructPointers works same as FillNilPointers,
// but fills only nil pointers to structs (other pointers will remain nil).
func FillNilStructPointers(s any) {
	opts := fillNilPointersOptions{
		onlyStructPtrs: true,
	}
	fillNilPointers(reflect.ValueOf(s), opts)
}

func fillNilPointers(v reflect.Value, opts fillNilPointersOptions) {
	// TODO: actually check CanAddr ?
	xmust.Eq(v.Kind(), reflect.Pointer)
	xmust.True(!v.IsNil())
	for v.Elem().Kind() == reflect.Pointer {
		if v.Elem().IsNil() {
			v.Elem().Set(reflect.New(v.Elem().Type().Elem()))
		}
		v = v.Elem()
	}
	xmust.Eq(v.Elem().Kind(), reflect.Struct, "input kind", xcheck.PrintValues)
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
