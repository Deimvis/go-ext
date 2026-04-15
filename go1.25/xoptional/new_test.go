package xoptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("empty/int", func(t *testing.T) {
		opt := New[int]()
		require.False(t, opt.HasValue())
		require.Nil(t, opt.ValuePtr())
		require.Panics(t, func() {
			opt.Value()
		})
		require.Panics(t, func() {
			*opt.ValuePtr() = 42
		})
	})
	t.Run("empty/string", func(t *testing.T) {
		opt := New[string]()
		require.False(t, opt.HasValue())
		require.Nil(t, opt.ValuePtr())
		require.Panics(t, func() {
			opt.Value()
		})
	})
	t.Run("empty/custom-type", func(t *testing.T) {
		type custom struct {
			Value string
		}
		opt := New[custom]()
		require.False(t, opt.HasValue())
		require.Nil(t, opt.ValuePtr())
		require.Panics(t, func() {
			opt.Value()
		})
	})
	t.Run("with-value", func(t *testing.T) {
		opt := New("hello world")
		require.True(t, opt.HasValue())
		require.NotNil(t, opt.ValuePtr())
		require.NotPanics(t, func() {
			opt.Value()
		})
	})
	t.Run("with-multiple-values", func(t *testing.T) {
		require.Panics(t, func() {
			New("value1", "value2")
		})
	})
}
