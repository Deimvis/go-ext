//go:build go1.25

package xwrapcallobs_test

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcallobs"
	"github.com/stretchr/testify/require"
)

func TestCallStackInd_go125(t *testing.T) {
	type C = context.Context

	t.Run("valid-active-stack-ind/no-mw", func(t *testing.T) {
		obs, view := xwrapcallobs.NewCallStackInd[C]()
		taskStarted := make(chan struct{})
		finishTask := make(chan struct{})
		fn := xwrapcall.New[C]().
			ObservingEvents(obs).
			Do(func(c C) error {
				close(taskStarted)
				<-finishTask
				return nil
			})
		_, ok := view.ActiveStackIndex(0)
		require.False(t, ok)

		synctest.Test(t, func(t *testing.T) {
			go func() {
				_ = fn(context.Background())
			}()
			<-taskStarted
			aInd, ok := view.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, xwrapcall.StackInd(0), aInd)

			close(finishTask)
			synctest.Wait()
			_, ok = view.ActiveStackIndex(0)
			require.False(t, ok)
		})
	})
	t.Run("valid-active-stack-ind/1-mw", func(t *testing.T) {
		obs, view := xwrapcallobs.NewCallStackInd[C]()
		mw1PreStarted := make(chan struct{})
		mw1PreFinish := make(chan struct{})
		mw1PostStarted := make(chan struct{})
		mw1PostFinish := make(chan struct{})
		taskStarted := make(chan struct{})
		finishTask := make(chan struct{})
		fn := xwrapcall.New[C]().
			ObservingEvents(obs).
			With(
				func(c C, next xwrapcall.Next[C]) (C, error) {
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
		_, ok := view.ActiveStackIndex(0)
		require.False(t, ok)

		synctest.Test(t, func(t *testing.T) {
			var aInd xwrapcall.StackInd

			go func() {
				_ = fn(context.Background())
			}()
			<-mw1PreStarted
			aInd, ok = view.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, xwrapcall.StackInd(0), aInd)

			close(mw1PreFinish)
			<-taskStarted
			aInd, ok = view.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, xwrapcall.StackInd(1), aInd)

			close(finishTask)
			<-mw1PostStarted
			aInd, ok = view.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, xwrapcall.StackInd(0), aInd)

			close(mw1PostFinish)
			synctest.Wait()
			_, ok = view.ActiveStackIndex(0)
			require.False(t, ok)
		})
	})
	t.Run("valid-active-stack-ind/second-run", func(t *testing.T) {
		obs, view := xwrapcallobs.NewCallStackInd[C]()
		taskStarted := make(chan struct{})
		finishTask := make(chan struct{})
		fn := xwrapcall.New[C]().
			ObservingEvents(obs).
			Do(func(c C) error {
				close(taskStarted)
				<-finishTask
				return nil
			})
		_, ok := view.ActiveStackIndex(0)
		require.False(t, ok)

		close(finishTask)
		_ = fn(context.Background())
		_, ok = view.ActiveStackIndex(0)
		require.False(t, ok)
		taskStarted = make(chan struct{})
		finishTask = make(chan struct{})

		synctest.Test(t, func(t *testing.T) {
			go func() {
				_ = fn(context.Background())
			}()
			<-taskStarted
			aInd, ok := view.ActiveStackIndex(0)
			require.True(t, ok)
			require.Equal(t, xwrapcall.StackInd(0), aInd)

			close(finishTask)
			synctest.Wait()
			_, ok = view.ActiveStackIndex(0)
			require.False(t, ok)
		})
	})
}
