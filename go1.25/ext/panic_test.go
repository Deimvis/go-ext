package ext

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOnPanic(t *testing.T) {
	{

		bug := true
		func() {
			defer func() {
				recover()
			}()
			defer OnPanic(func(any) {
				bug = false
			})
			panic("error")
		}()
		require.False(t, bug)
	}
	{

		bug := false
		func() {
			defer func() {
				recover()
			}()
			defer OnPanic(func(any) {
				bug = true
			})
			// no panic
		}()
		require.False(t, bug)
	}
	{

		bug := true
		func() {
			defer func() {
				recover()
			}()
			defer OnPanic(func(r any) {
				bug = r.(bool)
			})
			panic(false)
		}()
		require.False(t, bug)
	}
	// https://stackoverflow.com/a/49344592
	{

		bug := false
		func() {
			defer func() {
				recover()
			}()
			defer func(int) {
				OnPanic(func(any) {
					bug = true
				})
			}(42)
			panic("error")
		}()
		require.False(t, bug)
	}
}
