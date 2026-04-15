package xcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Run("new/empty", func(t *testing.T) {
		cm := NewMap(context.Background())
		var ok bool
		_, ok = cm.Get("a")
		require.False(t, ok)
		_, ok = cm.Get("b")
		require.False(t, ok)
		_, ok = cm.Get(1)
		require.False(t, ok)
	})
	t.Run("new/context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "value")
		cm := NewMap(ctx)
		require.Equal(t, "value", cm.Context().Value("key"))
	})
	t.Run("set", func(t *testing.T) {
		cm := NewMap(context.Background())
		var v any
		var ok bool
		_, ok = cm.Get("key")
		require.False(t, ok)

		cm.Set("key", "value")
		v, ok = cm.Get("key")
		require.True(t, ok)
		require.Equal(t, "value", v)

		cm.Set("key", "value2")
		v, ok = cm.Get("key")
		require.True(t, ok)
		require.Equal(t, "value2", v)
	})
	t.Run("context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "value")
		cm := NewMap(ctx)
		require.Equal(t, "value", cm.Context().Value("key"))

		cm.Set("key2", "value2")
		require.Equal(t, "value", cm.Context().Value("key"))
		require.Equal(t, "value2", cm.Context().Value("key2"))
	})
}
