package xreflect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndirect(t *testing.T) {
	for _, tc := range indirectTests {
		t.Run(tc.title, func(t *testing.T) {
			cur := tc.inp
			for _, exp := range tc.expUntilFalse {
				act, actOk := Indirect(cur)
				require.Equal(t, true, actOk)
				requireReflectEq(t, exp, act)
				cur = act
			}
			_, actOk := Indirect(cur)
			require.Equal(t, false, actOk)
		})
	}
}

func TestRecursiveIndirect(t *testing.T) {
	for _, tc := range indirectTests {
		t.Run(tc.title, func(t *testing.T) {
			act := RecursiveIndirect(tc.inp)
			exp := tc.inp
			if len(tc.expUntilFalse) > 0 {
				exp = tc.expUntilFalse[len(tc.expUntilFalse)-1]
			}
			requireReflectEq(t, exp, act)
		})
	}
}

func requireReflectEq(t *testing.T, exp reflect.Value, act reflect.Value) {
	require.Equal(t, exp.Kind(), act.Kind())
	if act.Kind() != reflect.Invalid {
		require.Equal(t, exp.Type(), act.Type())
		if act.IsZero() {
			require.True(t, exp.IsZero())
		} else {
			require.Equal(t, exp.Interface(), act.Interface())
		}
	}
}

var intV int = 42
var v = myType{V: 42}
var indirectTests = []indirectTc{
	{
		"int",
		reflect.ValueOf(intV),
		[]reflect.Value{},
	},
	{
		"*int/value",
		reflect.ValueOf(&intV),
		[]reflect.Value{
			reflect.ValueOf(intV),
		},
	},
	{
		"*int/nil",
		reflect.ValueOf((*int)(nil)),
		[]reflect.Value{},
	},
	{
		"**int/value",
		reflect.ValueOf(ptr(&intV)),
		[]reflect.Value{
			reflect.ValueOf(&intV),
			reflect.ValueOf(intV),
		},
	},
	{
		"**int/nil",
		reflect.ValueOf((**int)(nil)),
		[]reflect.Value{},
	},
	// Note that for interface tests we wrap top-level
	// interface (if any) with pointer and call Elem(),
	// because reflect.ValueOf automatically
	// indirects top-level interface from input value
	// (because top-level interface is converted to any,
	// and ValueOf can not distinguish this from case
	// where non-interface input value
	// was implicitly wrapped with interface any)
	{
		"interface(value)",
		reflect.ValueOf(ptr(myInterface(&v))).Elem(),
		[]reflect.Value{
			reflect.ValueOf(&v),
			reflect.ValueOf(v),
		},
	},
	{

		"interface(interface(*value))",
		reflect.ValueOf(ptr(myLimitedInterface(myInterface(&v)))).Elem(),
		[]reflect.Value{
			// second interface actually converts first one to second,
			// so after first indirection we immediately get underlying v
			reflect.ValueOf(&v),
			reflect.ValueOf(v),
		},
	},
	{
		"reflect.Invalid",
		reflect.ValueOf(v).MethodByName("nonexisting_method"),
		[]reflect.Value{},
	},
}

type indirectTc struct {
	title         string
	inp           reflect.Value
	expUntilFalse []reflect.Value
}

type myInterface interface {
	SetValue(int)
	Value() int
}

type myLimitedInterface interface {
	Value() int
}

type myType struct {
	V int
}

var _ myInterface = (*myType)(nil)
var _ myLimitedInterface = (*myType)(nil)

func (mt myType) Value() int {
	return mt.V
}

func (mt *myType) SetValue(v int) {
	mt.V = v
}

func ptr[T any](v T) *T {
	return &v
}
