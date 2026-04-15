//go:build go1.25

package xwrapcall

import (
	"context"
	"testing"
	"testing/synctest"

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

func TestNew_DebugInfo_go125(t *testing.T) {
	type C = context.Context

	t.Run("valid-active-stack-ind/no-mw", func(t *testing.T) {
		var d Debug = nil
		taskStarted := make(chan struct{})
		finishTask := make(chan struct{})
		fn := New[C]().
			ExportingDebugInfo(&d).
			Do(func(c C) error {
				close(taskStarted)
				<-finishTask
				return nil
			})
		require.NotNil(t, d)
		_, ok := d.ActiveStackIndex(0)
		require.False(t, ok)

		synctest.Test(t, func(t *testing.T) {
			go func() {
				_ = fn(context.Background())
			}()
			<-taskStarted
			aInd, ok := d.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, StackInd(0), aInd)

			close(finishTask)
			synctest.Wait()
			_, ok = d.ActiveStackIndex(0)
			require.False(t, ok)
		})
	})
	t.Run("valid-active-stack-ind/1-mw", func(t *testing.T) {
		var d Debug = nil
		mw1PreStarted := make(chan struct{})
		mw1PreFinish := make(chan struct{})
		mw1PostStarted := make(chan struct{})
		mw1PostFinish := make(chan struct{})
		taskStarted := make(chan struct{})
		finishTask := make(chan struct{})
		fn := New[C]().
			ExportingDebugInfo(&d).
			With(
				func(c C, next Next[C]) (C, error) {
					close(mw1PreStarted)
					<-mw1PreFinish
					c, err := next(c)
					close(mw1PostStarted)
					<-mw1PostFinish
					return c, err
				},
			).
			Do(func(c C) error {
				close(taskStarted)
				<-finishTask
				return nil
			})
		require.NotNil(t, d)
		_, ok := d.ActiveStackIndex(0)
		require.False(t, ok)

		synctest.Test(t, func(t *testing.T) {
			var aInd StackInd

			go func() {
				_ = fn(context.Background())
			}()
			<-mw1PreStarted
			aInd, ok = d.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, StackInd(0), aInd)

			close(mw1PreFinish)
			<-taskStarted
			aInd, ok = d.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, StackInd(1), aInd)

			close(finishTask)
			<-mw1PostStarted
			aInd, ok = d.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, StackInd(0), aInd)

			close(mw1PostFinish)
			synctest.Wait()
			_, ok = d.ActiveStackIndex(0)
			require.False(t, ok)
		})
	})
	t.Run("valid-active-stack-ind/second-run", func(t *testing.T) {
		var d Debug = nil
		taskStarted := make(chan struct{})
		finishTask := make(chan struct{})
		fn := New[C]().
			ExportingDebugInfo(&d).
			Do(func(c C) error {
				close(taskStarted)
				<-finishTask
				return nil
			})
		require.NotNil(t, d)
		_, ok := d.ActiveStackIndex(0)
		require.False(t, ok)

		close(finishTask)
		_ = fn(context.Background())
		_, ok = d.ActiveStackIndex(0)
		require.False(t, ok)
		taskStarted = make(chan struct{})
		finishTask = make(chan struct{})

		synctest.Test(t, func(t *testing.T) {
			go func() {
				_ = fn(context.Background())
			}()
			<-taskStarted
			aInd, ok := d.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, StackInd(0), aInd)

			close(finishTask)
			synctest.Wait()
			_, ok = d.ActiveStackIndex(0)
			require.False(t, ok)
		})
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
