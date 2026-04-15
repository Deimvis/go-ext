//go:build !debug

package xinvar

import (
	"testing"
)

func TestTrue(t *testing.T) {
	{
		// do nothing
		True(false)
	}
	{
		// do nothing
		True(true)
	}
}
