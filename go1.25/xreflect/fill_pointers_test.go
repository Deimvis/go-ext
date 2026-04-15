package xreflect_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xreflect"
)

func TestFillNilPointers(t *testing.T) {
	var zeroInt int
	var zeroStr string
	t.Run("no-depth", func(t *testing.T) {
		type T struct {
			I *int
			S *string
		}
		v := T{}
		xreflect.FillNilPointers(&v)
		require.NotNil(t, v.I)
		require.Equal(t, zeroInt, *v.I)
		require.NotNil(t, v.S)
		require.Equal(t, zeroStr, *v.S)
	})
	t.Run("no-depth-containers", func(t *testing.T) {
		type T struct {
			S []int
			M map[string]int
		}
		v := T{}
		xreflect.FillNilPointers(&v)
		require.NotNil(t, v.S)
		require.NotNil(t, v.M)
	})
	t.Run("unexported-fields", func(t *testing.T) {
		type T struct {
			I *int
			s *string
		}
		v := T{}
		xreflect.FillNilPointers(&v)
		require.NotNil(t, v.I)
		require.Equal(t, zeroInt, *v.I)
		require.Nil(t, v.s)
	})
	t.Run("depth-1", func(t *testing.T) {
		type T struct {
			U struct {
				I *int
			}
		}
		v := T{}
		xreflect.FillNilPointers(&v)
		require.NotNil(t, v.U)
		require.NotNil(t, v.U.I)
		require.Equal(t, zeroInt, *v.U.I)
	})
	t.Run("depth-1-but-intermediate-struct-filled", func(t *testing.T) {
		type T struct {
			U struct {
				I *int
			}
		}
		v := T{U: struct{ I *int }{}}
		xreflect.FillNilPointers(&v)
		require.NotNil(t, v.U)
		require.NotNil(t, v.U.I)
		require.Equal(t, zeroInt, *v.U.I)
	})
}

func TestFillNilStructPointers(t *testing.T) {
	t.Run("ignore-non-struct-pointers", func(t *testing.T) {
		type T struct {
			I *int
			S *string
		}
		v := T{}
		xreflect.FillNilStructPointers(&v)
		require.Nil(t, v.I)
		require.Nil(t, v.S)
	})
	t.Run("depth-1", func(t *testing.T) {
		type T struct {
			U struct {
				I *int
			}
			I *int
		}
		v := T{}
		xreflect.FillNilStructPointers(&v)
		require.Nil(t, v.I)
		require.NotNil(t, v.U)
		require.Nil(t, v.U.I)
	})
}

type fillNilPointersTestCase[T any] struct {
	orig T
	exp  T
}
