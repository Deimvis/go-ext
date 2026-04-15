package xoptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptional_Value(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		o := New("hello world")
		require.Equal(t, "hello world", o.Value())
		v := o.Value()
		v = "new string"
		_ = v
		require.Equal(t, "hello world", o.Value())
	})
	t.Run("*string", func(t *testing.T) {
		s := "hello world"
		o := New(&s)
		require.Equal(t, &s, o.Value())
		(*o.Value()) = "new string"
		require.Equal(t, "new string", *o.Value())
		require.Equal(t, "new string", s)
	})
	t.Run("custom-type", func(t *testing.T) {
		type user struct {
			Name string
		}
		o := New(user{Name: "dmitriy"})
		require.Equal(t, "dmitriy", o.Value().Name)
		v := o.Value()
		v.Name = "aleksandr"
		require.Equal(t, "dmitriy", o.Value().Name)
	})
	t.Run("*custom-type", func(t *testing.T) {
		type user struct {
			Name string
		}
		u := user{Name: "dmitriy"}
		o := New(&u)
		require.Equal(t, &u, o.Value())
		o.Value().Name = "aleksandr"
		require.Equal(t, "aleksandr", o.Value().Name)
		require.Equal(t, "aleksandr", u.Name)
	})
}

func TestOptional_ValuePtr(t *testing.T) {
	t.Run("mutable/string", func(t *testing.T) {
		o := New("hello world")
		require.Equal(t, "hello world", o.Value())
		*o.ValuePtr() = "new string"
		require.Equal(t, "new string", o.Value())
	})
	t.Run("mutable/custom-type", func(t *testing.T) {
		type user struct {
			Name string
		}
		o := New(user{Name: "dmitriy"})
		require.Equal(t, "dmitriy", o.Value().Name)
		o.ValuePtr().Name = "aleksandr"
		require.Equal(t, "aleksandr", o.Value().Name)
	})
}

func TestOptional_SetValue(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		o := New("hello world")
		require.Equal(t, "hello world", o.Value())
		o.SetValue("new string")
		require.Equal(t, "new string", o.Value())
	})
}

func TestOptional_Reset(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		o := New("hello world")
		require.True(t, o.HasValue())
		require.Equal(t, "hello world", o.Value())
		o.Reset()
		require.False(t, o.HasValue())
	})
}
