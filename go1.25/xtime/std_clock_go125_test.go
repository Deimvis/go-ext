//go:build go1.25

package xtime

import (
	"runtime"
	"strconv"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

func TestStdClock_go125(t *testing.T) {
	t.Run("Now", func(t *testing.T) {
		// if goVersion()[1] < 25 {
		// 	t.Skip()
		// }
		synctest.Test(t, func(t *testing.T) {
			c := NewStdClock()
			exp := time.Now().UnixNano()
			act := c.Now().UnixNano()
			// in bubble time should be exactly equal
			require.Equal(t, exp, act)
		})
	})
}

func goVersion() [3]int {
	version := runtime.Version()
	version = strings.TrimPrefix(version, "go")
	parts := strings.Split(version, ".")
	xmust.Eq(len(parts), 3)
	res := [3]int{}
	for i, p := range parts {
		res[i] = xmust.Do(strconv.Atoi(p))
	}
	return res
}
