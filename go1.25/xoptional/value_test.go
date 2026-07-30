package xoptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromValueWithOk(t *testing.T) {
	t.Run("ok=true/holds-value", func(t *testing.T) {
		o := FromValueWithOk("hello", true)
		require.True(t, o.HasValue())
		require.Equal(t, "hello", o.Value())
	})
	t.Run("ok=false/empty", func(t *testing.T) {
		o := FromValueWithOk("hello", false)
		require.False(t, o.HasValue())
	})
	t.Run("ok=true/zero-value-is-still-set", func(t *testing.T) {
		o := FromValueWithOk("", true)
		require.True(t, o.HasValue())
		require.Equal(t, "", o.Value())
	})
	t.Run("round-trip-with-ValueWithOk", func(t *testing.T) {
		exp := 42
		o := FromValueWithOk(exp, true)
		act, ok := ValueWithOk(o)
		require.True(t, ok)
		require.Equal(t, exp, act)
	})
	t.Run("round-trip-empty", func(t *testing.T) {
		o := FromValueWithOk(42, false)
		_, ok := ValueWithOk(o)
		require.False(t, ok)
	})
}
