package xwrapcallobs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcallobs"
	"github.com/stretchr/testify/require"
)

func TestIssues(t *testing.T) {
	type C = context.Context
	passThrough := func(c C, next xwrapcall.Next[C]) (C, error) { return next(c) }

	t.Run("empty-before-run", func(t *testing.T) {
		_, view := xwrapcallobs.NewIssues[C]()
		require.Empty(t, view.Events())
	})
	t.Run("normal-return/no-events", func(t *testing.T) {
		obs, view := xwrapcallobs.NewIssues[C]()
		fn := xwrapcall.New[C]().
			With(passThrough, passThrough).
			ObservingEvents(obs).
			Do(func(c C) error { return nil })
		require.NoError(t, fn(context.Background()))
		require.Empty(t, view.Events())
	})
	t.Run("error-return/records-frames-that-carry-error", func(t *testing.T) {
		obs, view := xwrapcallobs.NewIssues[C]()
		myErr := errors.New("boom")
		fn := xwrapcall.New[C]().
			With(passThrough, passThrough).
			ObservingEvents(obs).
			Do(func(c C) error { return myErr })
		err := fn(context.Background())
		require.ErrorIs(t, err, myErr)
		events := view.Events()
		require.NotEmpty(t, events)
		for _, e := range events {
			require.ErrorIs(t, e.Err, myErr)
			require.False(t, e.Panicked)
		}
	})
	t.Run("panic/marks-all-unwind-frames", func(t *testing.T) {
		obs, view := xwrapcallobs.NewIssues[C]()
		panicking := func(c C, next xwrapcall.Next[C]) (C, error) {
			panic("boom")
		}
		fn := xwrapcall.New[C]().
			With(passThrough, panicking).
			ObservingEvents(obs).
			Do(func(c C) error { return nil })
		require.Panics(t, func() {
			_ = fn(context.Background())
		})
		events := view.Events()
		require.Len(t, events, 2)
		// Leave events fire innermost-first: index 1 (panicking) then index 0 (passThrough)
		require.Equal(t, xwrapcall.StackInd(1), events[0].StackIndex)
		require.True(t, events[0].Panicked)
		require.Equal(t, xwrapcall.StackInd(0), events[1].StackIndex)
		require.True(t, events[1].Panicked)
	})
	t.Run("events-returns-snapshot", func(t *testing.T) {
		obs, view := xwrapcallobs.NewIssues[C]()
		myErr := errors.New("boom")
		fn := xwrapcall.New[C]().
			ObservingEvents(obs).
			Do(func(c C) error { return myErr })
		_ = fn(context.Background())

		e1 := view.Events()
		require.NotEmpty(t, e1)
		e1[0].StackIndex = 999
		e2 := view.Events()
		require.Equal(t, xwrapcall.StackInd(0), e2[0].StackIndex)
	})
}
