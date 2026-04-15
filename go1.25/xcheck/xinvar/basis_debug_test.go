//go:build debug

package xinvar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrue(t *testing.T) {
	exited := false
	exitFn = func(string) { exited = true }
	{
		exited = false
		True(false)
		require.True(t, exited)
	}
	{
		exited = false
		True(true)
		require.False(t, exited)
	}
}
