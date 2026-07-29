//go:build go1.25

package xwrapcall

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_RuntimeInfo(t *testing.T) {
	t.Run("valid-caller-stack-ind", func(t *testing.T) {
		type C = *cctx
		requireInd := func(t *testing.T, c C, i StackInd) {
			require.Equal(t, i, MustCallerRuntime(c).CallerStackIndex())
		}
		fn := New[C]().
			With(
				func(c C, next Next[C]) (C, error) {
					requireInd(t, c, 0)
					c, err := next(c)
					requireInd(t, c, 0)
					return c, err
				},
				func(c C, next Next[C]) (C, error) {
					requireInd(t, c, 1)
					c, err := next(c)
					requireInd(t, c, 1)
					return c, err
				},
			).
			Do(func(c C) error {
				requireInd(t, c, 2)
				return nil
			})

		c := &cctx{Context: context.Background()}
		err := fn(c)
		require.NoError(t, err)
	})
}

type cctx struct {
	context.Context
}

var _ CopyableContext[*cctx] = (*cctx)(nil)

func (c *cctx) StdContext() context.Context {
	return c.Context
}

func (c *cctx) CopyOnto(ctx context.Context) *cctx {
	cCopy := &cctx{
		Context: ctx,
	}
	return cCopy
}
