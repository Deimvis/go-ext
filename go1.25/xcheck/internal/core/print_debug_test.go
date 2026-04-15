//go:build debug

package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMsg(t *testing.T) {
	t.Run("no-args/ok", func(t *testing.T) {
		_validateMsg("valid")
	})
	t.Run("no-args/not-string-first-fails", func(t *testing.T) {
		require.Panics(t, func() {
			_validateMsg(42)
		})
	})
	t.Run("args/ok", func(t *testing.T) {
		_validateMsg("mymsg: %d, %s", 42, "123")
	})
	t.Run("args/invalid-args-number-fails/more", func(t *testing.T) {
		require.Panics(t, func() {
			_validateMsg("mymsg: %d, %s", 42, "123", "whoooa")
		})
	})
	t.Run("args/invalid-args-number-fails/less", func(t *testing.T) {
		require.Panics(t, func() {
			_validateMsg("mymsg: %d, %s", 42)
		})
	})
	t.Run("args/invalid-arg-type-fails/less", func(t *testing.T) {
		require.Panics(t, func() {
			_validateMsg("mymsg: %d, %s", 42, 123)
		})
	})
}
