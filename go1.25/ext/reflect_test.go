package ext

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnwrapStruct(t *testing.T) {
	testCases := []struct {
		title    string
		obj      interface{}
		expected []interface{}
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"plain struct",
			R_A{Key: 42, Value: "hello"},
			[]interface{}{
				42,
				"hello",
			},
		},
		{
			"pointer to plain struct",
			&R_A{Key: 42, Value: "hello"},
			[]interface{}{
				42,
				"hello",
			},
		},
		{
			"embedded struct",
			R_B{R_A: R_A{Key: 42, Value: "hello"}, Other: 0},
			[]interface{}{
				42,
				"hello",
				0,
			},
		},
		{
			"exported struct",
			R_C{A: R_A{Key: 42, Value: "hello"}, Other: 0},
			[]interface{}{
				R_A{Key: 42, Value: "hello"},
				0,
			},
		},
		{
			"unexported struct",
			R_D{a: R_A{Key: 42, Value: "hello"}, Other: 0},
			[]interface{}{
				0,
			},
		},
		{
			"recursively embedded struct",
			R_E{R_B: R_B{R_A: R_A{Key: 42, Value: "hello"}, Other: 0}, Other2: 1},
			[]interface{}{
				42,
				"hello",
				0,
				1,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := UnwrapStruct(tc.obj)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetBaseFnName(t *testing.T) {
	testCases := []struct {
		title    string
		fn       interface{}
		expected string
	}{
		{
			"foo",
			foo,
			"foo",
		},
		{
			"bar",
			bar,
			"bar",
		},
		{
			"test func (TestUnwrapStruct)",
			TestUnwrapStruct,
			"TestUnwrapStruct",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%d", tc.title, i), func(t *testing.T) {
			actual := GetBaseFnName(tc.fn)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetFieldAddr(t *testing.T) {
	a := R_A{Key: 42}
	actual, err := GetFieldAddr(&a, "Key")
	require.NoError(t, err)
	require.Equal(t, &a.Key, actual)
	actualTyped := actual.(*int)
	require.NotNil(t, actualTyped)
	*actualTyped = 55
	require.Equal(t, 55, a.Key)

	b := R_B{R_A: a}
	actual, err = GetFieldAddr(&b, "R_A")
	require.NoError(t, err)
	require.Equal(t, &b.R_A, actual)
}

func TestFillNilSlices(t *testing.T) {
	require.Panics(t, func() { FillNilSlices(42) })
	require.Panics(t, func() { FillNilSlices("hello") })
	require.Panics(t, func() { FillNilSlices(nil) })
	require.Panics(t, func() { FillNilSlices(R_A{}) })
	{
		x := struct {
			Key int
		}{
			Key: 42,
		}
		xcopy := x
		FillNilSlices(&x)
		require.Equal(t, xcopy, x)
	}
	{
		x := struct {
			Key int
			Arr []int
		}{
			Key: 42,
			Arr: nil,
		}
		xcopy := x
		FillNilSlices(&x)
		require.NotEqual(t, xcopy, x)
		require.NotEqual(t, nil, x.Arr)
		require.NotEqual(t, []int(nil), x.Arr)
		require.Equal(t, []int{}, x.Arr)
	}
	{
		x := struct {
			Key         int
			notExported []int
		}{
			Key:         42,
			notExported: nil,
		}
		FillNilSlices(&x)
		require.Equal(t, []int(nil), x.notExported)
	}
	{
		var x myInterface
		x = &myImpl{
			Arr: nil,
		}
		FillNilSlices(&x)
		require.Equal(t, []int{}, x.(*myImpl).Arr)
	}
	{
		type subT struct {
			Arr []int
		}
		x := struct {
			Sub subT
		}{
			Sub: subT{
				Arr: nil,
			},
		}
		FillNilSlices(&x)
		require.Equal(t, []int{}, x.Sub.Arr)
	}
	{
		type elemT struct {
			SubArr []int
		}
		type subT struct {
			Arr []elemT
		}
		x := struct {
			Sub subT
		}{
			Sub: subT{
				Arr: []elemT{
					{
						SubArr: nil,
					},
				},
			},
		}
		FillNilSlices(&x)
		require.Equal(t, []int{}, x.Sub.Arr[0].SubArr)

	}
}

func TestStructTypesAreEquivalent(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		act := StructTypesAreEquivalent(Equiv1{}, Equiv2{})
		require.Equal(t, true, act)
	})
	t.Run("not_equivalent_1", func(t *testing.T) {
		act := StructTypesAreEquivalent(struct{}{}, Equiv2{})
		require.Equal(t, false, act)
	})
	t.Run("not_equivalent_2", func(t *testing.T) {
		act := StructTypesAreEquivalent(Equiv1{}, struct{}{})
		require.Equal(t, false, act)
	})
	t.Run("equivalent_filled_objects", func(t *testing.T) {
		x := 42
		obj1 := Equiv1{
			Int:      42,
			String:   "mystring",
			Bool:     true,
			Byte:     12,
			IntPtr:   &x,
			IntSlice: []int{1, 2, 3},
			SubStruct: struct{ Value int64 }{
				Value: 999,
			},
		}
		y := 24
		obj2 := Equiv2{
			Int:      24,
			String:   "notmystring",
			Bool:     false,
			Byte:     21,
			IntPtr:   &y,
			IntSlice: []int{3, 2, 1},
			SubStruct: struct{ Value int64 }{
				Value: 666,
			},
		}
		act := StructTypesAreEquivalent(obj1, obj2)
		require.Equal(t, true, act)
	})
}

func TestConvertEquivalentStructs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		x := 42
		obj1 := Equiv1{
			Int:      42,
			String:   "mystring",
			Bool:     true,
			Byte:     12,
			IntPtr:   &x,
			IntSlice: []int{1, 2, 3},
			SubStruct: struct{ Value int64 }{
				Value: 999,
			},
		}
		var obj2 Equiv2
		obj1copy := obj1
		err := ConvertEquivalentStructs(obj1, &obj2)
		require.NoError(t, err)
		testEqual11(t, obj1copy, obj1)
		testEqual12(t, obj1, obj2)
	})
	t.Run("filled_out_object", func(t *testing.T) {
		x := 42
		obj1 := Equiv1{
			Int:      42,
			String:   "mystring",
			Bool:     true,
			Byte:     12,
			IntPtr:   &x,
			IntSlice: []int{1, 2, 3},
			SubStruct: struct{ Value int64 }{
				Value: 999,
			},
		}
		y := 24
		obj2 := Equiv2{
			Int:      24,
			String:   "notmystring",
			Bool:     false,
			Byte:     21,
			IntPtr:   &y,
			IntSlice: []int{3, 2, 1},
			SubStruct: struct{ Value int64 }{
				Value: 666,
			},
		}
		obj1copy := obj1
		err := ConvertEquivalentStructs(obj1, &obj2)
		require.NoError(t, err)
		testEqual11(t, obj1copy, obj1)
		testEqual12(t, obj1, obj2)
	})
}

type R_A struct {
	Key   int
	Value string
}

type R_B struct {
	R_A
	Other int
}

type R_C struct {
	A     R_A
	Other int
}

type R_D struct {
	a     R_A
	Other int
}

type R_E struct {
	R_B
	Other2 int
}

func foo()          {}
func bar(int) error { return nil }

type getFieldAddrTestCase struct {
	title         string
	s             interface{}
	fieldName     string
	expectedValue interface{}
	expectedError bool
}

type myInterface interface {
	Foo() int
}

type myImpl struct {
	Arr []int
}

func (*myImpl) Foo() int {
	return 42
}

type Equiv1 struct {
	Int       int
	String    string
	Bool      bool
	Byte      byte
	IntPtr    *int
	IntSlice  []int
	SubStruct struct {
		Value int64
	}
}

type Equiv2 struct {
	IntPtr    *int
	Byte      byte
	Int       int
	Bool      bool
	SubStruct struct {
		Value int64
	}
	IntSlice []int
	String   string
}

func testEqual11(t *testing.T, e1 Equiv1, e2 Equiv1) {
	require.Equal(t, e1.Int, e2.Int)
	require.Equal(t, e1.String, e2.String)
	require.Equal(t, e1.Bool, e2.Bool)
	require.Equal(t, e1.Byte, e2.Byte)
	require.Equal(t, e1.IntPtr, e2.IntPtr)
	require.Equal(t, e1.IntSlice, e2.IntSlice)
	require.Equal(t, e1.SubStruct.Value, e2.SubStruct.Value)
}

func testEqual12(t *testing.T, e1 Equiv1, e2 Equiv2) {
	require.Equal(t, e1.Int, e2.Int)
	require.Equal(t, e1.String, e2.String)
	require.Equal(t, e1.Bool, e2.Bool)
	require.Equal(t, e1.Byte, e2.Byte)
	require.Equal(t, e1.IntPtr, e2.IntPtr)
	require.Equal(t, e1.IntSlice, e2.IntSlice)
	require.Equal(t, e1.SubStruct.Value, e2.SubStruct.Value)
}
