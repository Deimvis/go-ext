package xoptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromPtr(t *testing.T) {
	t.Run("nil/int", func(t *testing.T) {
		var x *int = nil
		xopt := FromPtr(x)
		require.False(t, xopt.HasValue())
	})
	t.Run("value/int", func(t *testing.T) {
		var v int = 42
		x := &v
		xopt := FromPtr(x)
		require.True(t, xopt.HasValue())
		require.Equal(t, 42, xopt.Value())
		v = 123
		require.Equal(t, 42, xopt.Value())
	})
}

func TestToPtr(t *testing.T) {
	t.Run("nil/int", func(t *testing.T) {
		xopt := New[int]()
		x := ToPtr(xopt)
		require.Nil(t, x)
	})
	t.Run("value/int", func(t *testing.T) {
		xopt := New[int](42)
		x := ToPtr(xopt)
		require.NotNil(t, x)
		require.Equal(t, 42, *x)
		xopt.SetValue(123)
		require.Equal(t, 42, *x)
	})
}

func TestToPtrWithOk(t *testing.T) {
	t.Run("nil/int", func(t *testing.T) {
		xopt := New[int]()
		x, ok := ToPtrWithOk(xopt)
		require.False(t, ok)
		require.Nil(t, x)
	})
	t.Run("value/int", func(t *testing.T) {
		xopt := New[int](42)
		x, ok := ToPtrWithOk(xopt)
		require.True(t, ok)
		require.NotNil(t, x)
		require.Equal(t, 42, *x)
		xopt.SetValue(123)
		require.Equal(t, 42, *x)
	})
}
