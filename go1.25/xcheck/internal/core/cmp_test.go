package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotNilInterface(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		var x any = 42
		ok, _ := NotNilInterface(x)
		require.True(t, ok)
	})
	t.Run("fail", func(t *testing.T) {
		var x any
		ok, _ := NotNilInterface(x)
		require.False(t, ok)
	})
	t.Run("pass/nil-value", func(t *testing.T) {
		var x any = (*int)(nil)
		ok, _ := NotNilInterface(x)
		require.True(t, ok)
	})
}
