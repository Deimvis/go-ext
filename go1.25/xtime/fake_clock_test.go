//go:build debug

package xtime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakeClock(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		c := NewFakeClock()
		require.NotNil(t, c)
	})
}
