//go:build go1.25 && debug

package xtime

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFakeClock_go125(t *testing.T) {
	t.Run("Now/default", func(t *testing.T) {
		// if goVersion()[1] < 25 {
		// 	t.Skip()
		// }
		synctest.Test(t, func(t *testing.T) {
			c := NewFakeClock()
			exp := time.Now().UnixNano()
			act := c.Now().UnixNano()
			// in bubble time should be exactly equal
			require.Equal(t, exp, act)
		})
	})
	t.Run("Shift", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := NewFakeClock()
			t1 := c.Now()
			c.Shift(time.Minute)
			t2 := c.Now()
			require.Equal(t, t1.Add(time.Minute), t2)
			// persists
			t3 := c.Now()
			require.Equal(t, t3, t2)
		})
	})
	t.Run("Stop", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := NewFakeClock()
			t1 := c.Now()
			c.Stop()
			t2 := c.Now()
			require.Equal(t, t1, t2)
			// ignores shift
			c.Shift(time.Minute)
			t3 := c.Now()
			require.Equal(t, t3, t2)
		})
	})
	t.Run("StopAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := NewFakeClock()
			t1 := c.Now()
			c.StopAt(t1.Add(-time.Minute))
			t2 := c.Now()
			require.Equal(t, t1.Add(-time.Minute), t2)
			// persists
			t3 := c.Now()
			require.Equal(t, t3, t2)
		})
	})
	t.Run("Reset/after-shift", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := NewFakeClock()
			t1 := c.Now()
			c.Shift(time.Minute)
			t2 := c.Now()
			require.Equal(t, t1.Add(time.Minute), t2)
			c.Reset()
			t3 := c.Now()
			require.Equal(t, t1, t3)
		})
	})
	t.Run("Reset/after-stop", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := NewFakeClock()
			t1 := c.Now()
			c.Stop()
			c.Reset()
			t2 := c.Now()
			require.Equal(t, t1, t2)
		})
	})
}
