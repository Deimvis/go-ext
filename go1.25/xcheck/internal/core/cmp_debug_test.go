//go:build debug

package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotNilInterface_debug(t *testing.T) {
	t.Run("fatal/ptr", func(t *testing.T) {
		crashed, cancel := withFatalInterception()
		defer cancel()

		var ptr *int = nil

		_, _ = NotNilInterface(ptr)
		require.False(t, *crashed)
	})
	t.Run("fatal/slice", func(t *testing.T) {
		crashed, cancel := withFatalInterception()
		defer cancel()

		var s []int = nil

		_, _ = NotNilInterface(s)
		require.False(t, *crashed)
	})
	t.Run("fatal/map", func(t *testing.T) {
		crashed, cancel := withFatalInterception()
		defer cancel()

		var m map[int]int = nil

		_, _ = NotNilInterface(m)
		require.False(t, *crashed)
	})
	t.Run("fatal/chan", func(t *testing.T) {
		crashed, cancel := withFatalInterception()
		defer cancel()

		var c chan int = nil

		_, _ = NotNilInterface(c)
		require.False(t, *crashed)
	})
	t.Run("fatal/func", func(t *testing.T) {
		crashed, cancel := withFatalInterception()
		defer cancel()

		var f func() = nil

		_, _ = NotNilInterface(f)
		require.False(t, *crashed)
	})
}

func withFatalInterception() (*bool, func()) {
	crashed := false
	old := fatal
	new := func(format string, v ...any) {
		crashed = true
	}

	fatal = new
	cancel := func() {
		fatal = old
	}
	return &crashed, cancel
}
