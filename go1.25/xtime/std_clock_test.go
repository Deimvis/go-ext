package xtime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStdClock(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		c := NewStdClock()
		require.NotNil(t, c)
	})
	t.Run("not-configurable", func(t *testing.T) {
		c := NewStdClock()
		_, ok := c.(ConfigurableClock)
		require.False(t, ok)
	})
}
