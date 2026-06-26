package xmust

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoErr(t *testing.T) {
	t.Run("preserves-original-typed-error", func(t *testing.T) {
		orig := customErr{tag: "abc"}

		r := recoverPanic(func() { NoErr(orig) })
		require.NotNil(t, r, "NoErr must panic on non-nil error")
		rErr, ok := r.(error)
		require.True(t, ok, "panic value must be an error, got %T", r)

		var got customErr
		require.True(t, errors.As(rErr, &got), "errors.As must recover the original typed error from the panic value")
		require.Equal(t, "abc", got.tag)
	})

	t.Run("preserves-wrapped-error-chain", func(t *testing.T) {
		inner := customErr{tag: "inner"}
		wrapped := fmt.Errorf("outer: %w", inner)

		r := recoverPanic(func() { NoErr(wrapped) })
		require.NotNil(t, r)
		rErr, ok := r.(error)
		require.True(t, ok)

		var got customErr
		require.True(t, errors.As(rErr, &got))
		require.Equal(t, "inner", got.tag)
	})

	t.Run("preserves-sentinel-error", func(t *testing.T) {
		sentinel := errors.New("sentinel")

		r := recoverPanic(func() { NoErr(sentinel) })
		require.NotNil(t, r)
		rErr, ok := r.(error)
		require.True(t, ok)
		require.True(t, errors.Is(rErr, sentinel), "errors.Is must match the original sentinel error")
	})

	t.Run("does-not-panic-on-nil", func(t *testing.T) {
		require.NotPanics(t, func() { NoErr(nil) })
	})
}

func recoverPanic(fn func()) (r any) {
	defer func() { r = recover() }()
	fn()
	return nil
}

type customErr struct {
	tag string
}

func (e customErr) Error() string {
	return "custom: " + e.tag
}
