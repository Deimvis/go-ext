package xruntime

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXCaller(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cp, err := XCaller(0)
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(cp.File(), "ext/xruntime/extern_test.go"))
		require.Equal(t, 12, cp.Line())
		require.True(t, strings.HasSuffix(cp.PackagePath(), "ext/xruntime"))
	})
	t.Run("nested", func(t *testing.T) {
		foo := func() {
			cp, err := XCaller(0)
			require.NoError(t, err)
			require.True(t, strings.HasSuffix(cp.File(), "ext/xruntime/extern_test.go"))
			require.Equal(t, 20, cp.Line())
			require.True(t, strings.HasSuffix(cp.PackagePath(), "ext/xruntime"))
		}
		foo()
	})
}
