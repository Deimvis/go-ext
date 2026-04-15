package ext

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/Deimvis/go-ext/go1.25/xutf8string"
)

// UnwrapStruct returns slice of values of each field.
// It unwraps embedded fields and returns its subfields
// ignoring embedded fields themselves.
// It ignores unexported (lowercase) fields.
// The order of returned struct field values is guaranteed
// to match the order in which fields were declared in the struct definition.
func UnwrapStruct(s interface{}) []interface{} {
	if s == nil {
		return nil
	}
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	visibleFields := reflect.VisibleFields(v.Type())
	var fields []interface{}
	for _, f := range visibleFields {
		if f.Anonymous || !f.IsExported() {
			continue
		}
		fields = append(fields, v.FieldByIndex(f.Index).Interface())
	}
	return fields
}

// getStructFields returns slice of values of each field.
// It returns only fields that are declared in struct definition
// which means it ignores all embedded subfields.
// The order of returned struct field values is guaranteed
// to match the order in which fields were declared in the struct definition.
func getStructDefinitionFields(s interface{}) []interface{} {
	if s == nil {
		return nil
	}
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	var fields []interface{}
	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, v.Field(i).Interface())
	}
	return fields
}

func GetBaseFnName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

func GetField(s interface{}, fieldName string) (any, error) {
	if len(fieldName) == 0 {
		return nil, errors.New("no field name was given")
	}
	if !xutf8string.IsCapitalized(fieldName) {
		return nil, fmt.Errorf("can't access not exported field (one starting with lowercase letter) - %s", fieldName)
	}

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't get field of non-struct object (kind = %v)", v.Kind())
	}

	fieldV := v.FieldByName(fieldName)
	if !fieldV.IsValid() {
		return nil, nil
	}
	return fieldV.Interface(), nil
}

// GetFieldAddr returns a pointer to field of the given struct.
// Returns nil if no such field exists.
func GetFieldAddr(s interface{}, fieldName string) (any, error) {
	if len(fieldName) == 0 {
		return nil, errors.New("no field name was given")
	}
	if !xutf8string.IsCapitalized(fieldName) {
		return nil, fmt.Errorf("can't access not exported field (one starting with lowercase letter) - %s", fieldName)
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("can't get field address of non-pointer object (kind = %v)", v.Kind())
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't get field address of non-struct object (kind = %v)", v.Kind())
	}

	fieldV := v.FieldByName(fieldName)
	if !fieldV.IsValid() {
		return nil, nil
	}
	if !fieldV.CanAddr() {
		return nil, errors.New("field can't be addressed")
	}
	return fieldV.Addr().Interface(), nil
}

// TODO: migrate to xreflect.FillNilPointers
func FillNilSlices(obj interface{}) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer {
		panic("can't fill slices on non-pointer object")
	}
	v = extractInternalValue(v)
	if v.Kind() != reflect.Struct {
		panic("can't fill slices on non-struct object")
	}
	visibleFields := reflect.VisibleFields(v.Type())
	for _, f := range visibleFields {
		if f.Anonymous || !f.IsExported() {
			continue
		}
		field := v.FieldByIndex(f.Index)
		switch field.Kind() {
		case reflect.Struct:
			if field.CanAddr() {
				FillNilSlices(field.Addr().Interface())
			}
		case reflect.Slice:
			if field.IsNil() {
				newSlice := reflect.MakeSlice(field.Type(), 0, 0)
				field.Set(newSlice)
			}
			for i := 0; i < field.Len(); i++ {
				if field.Index(i).Kind() != reflect.Struct {
					break
				}
				if field.Index(i).CanAddr() {
					FillNilSlices(field.Index(i).Addr().Interface())
				}
			}
		}
	}
}

// StructTypesAreEquivalent checks whethet 2 struct types
// have the same fields (same names, same types, order insensetive).
func StructTypesAreEquivalent(obj1, obj2 interface{}) bool {
	type1 := reflect.TypeOf(obj1)
	type2 := reflect.TypeOf(obj2)

	for type1.Kind() == reflect.Pointer {
		type1 = type1.Elem()
	}
	for type2.Kind() == reflect.Pointer {
		type2 = type2.Elem()
	}

	if type1.Kind() != reflect.Struct || type2.Kind() != reflect.Struct {
		return false
	}

	if type1.NumField() != type2.NumField() {
		return false
	}

	name2type := make(map[string]reflect.Type)
	for i := 0; i < type1.NumField(); i++ {
		name2type[type1.Field(i).Name] = type1.Field(i).Type
	}
	for i := 0; i < type2.NumField(); i++ {
		if _, ok := name2type[type2.Field(i).Name]; !ok {
			return false
		}
		if type2.Field(i).Type != name2type[type2.Field(i).Name] {
			return false
		}
	}

	return true
}

// ConvertEquivalentStructs runs StructTypesAreEquivalent check
// and then fills objOut using objIn.
//
// var obj1 lib1.User = lib1.User{name: "john"}
// var obj2 lib2.User
// ConvertEquivalentStructs(obj1, &obj2)
func ConvertEquivalentStructs(objIn, objOut interface{}) error {
	vIn := reflect.ValueOf(objIn)
	if vIn.Kind() != reflect.Struct {
		return fmt.Errorf("conversion input object should be a struct, not an element of kind = %v", vIn.Kind())
	}
	vOut := reflect.ValueOf(objOut)
	if vOut.Kind() != reflect.Pointer {
		return fmt.Errorf("conversion output object should be a struct pointer, not an element of kind = %v", vOut.Kind())
	}
	if vOut.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("conversion output object should be a struct pointer, not a pointer to kind = %v", vOut.Kind())
	}
	if !StructTypesAreEquivalent(objIn, objOut) {
		return errors.New("struct types should be equivalent")
	}
	return mapstructure.Decode(objIn, objOut)
}

// https://github.com/go-playground/validator/blob/a947377040f8ebaee09f20d09a745ec369396793/util.go#L15
func extractInternalValue(current reflect.Value) reflect.Value {

BEGIN:
	switch current.Kind() {
	case reflect.Pointer:

		if current.IsNil() {
			return current
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:

		if current.IsNil() {
			return current
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current

	default:
		return current
	}
}
