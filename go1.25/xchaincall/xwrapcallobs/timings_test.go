package xwrapcallobs_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcallobs"
	"github.com/Deimvis/go-ext/go1.25/xtime"
	"github.com/stretchr/testify/require"
)

func TestTimings(t *testing.T) {
	t.Run("empty-before-run", func(t *testing.T) {
		_, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		require.Len(t, view.Events(), 0)
	})
	t.Run("no-mws/records-enter-and-leave", func(t *testing.T) {
		obs, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		fn := xwrapcall.New[context.Context]().
			ObservingEvents(obs).
			Do(func(c context.Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)

		events := view.Events()
		require.Len(t, events, 2)

		require.Equal(t, xwrapcallobs.FrameEnter, events[0].Kind)
		require.Equal(t, xwrapcall.StackInd(0), events[0].StackIndex)
		require.False(t, events[0].Ts.IsZero())

		require.Equal(t, xwrapcallobs.FrameLeave, events[1].Kind)
		require.Equal(t, xwrapcall.StackInd(0), events[1].StackIndex)
		require.False(t, events[1].Ts.IsZero())
	})
	t.Run("mws/full-trace/preserves-order", func(t *testing.T) {
		obs, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		passThrough := func(c context.Context, next xwrapcall.Next[context.Context]) (context.Context, error) {
			return next(c)
		}
		fn := xwrapcall.New[context.Context]().
			With(passThrough, passThrough).
			ObservingEvents(obs).
			Do(func(c context.Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)

		events := view.Events()
		require.Len(t, events, 6)

		expected := []struct {
			kind xwrapcallobs.TimingsEventKind
			ind  xwrapcall.StackInd
		}{
			{xwrapcallobs.FrameEnter, 0},
			{xwrapcallobs.FrameEnter, 1},
			{xwrapcallobs.FrameEnter, 2},
			{xwrapcallobs.FrameLeave, 2},
			{xwrapcallobs.FrameLeave, 1},
			{xwrapcallobs.FrameLeave, 0},
		}
		for i, exp := range expected {
			require.Equal(t, exp.kind, events[i].Kind, "event %d", i)
			require.Equal(t, exp.ind, events[i].StackIndex, "event %d", i)
		}
	})
	t.Run("timestamps-monotonic-non-decreasing", func(t *testing.T) {
		obs, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		slowMw := func(c context.Context, next xwrapcall.Next[context.Context]) (context.Context, error) {
			time.Sleep(time.Microsecond)
			return next(c)
		}
		fn := xwrapcall.New[context.Context]().
			With(slowMw, slowMw).
			ObservingEvents(obs).
			Do(func(c context.Context) error {
				time.Sleep(time.Microsecond)
				return nil
			})
		err := fn(context.Background())
		require.NoError(t, err)

		events := view.Events()
		require.Len(t, events, 6)
		for i := 1; i < len(events); i++ {
			require.True(t, events[i].Ts.After(events[i-1].Ts),
				"event %d timestamp went backwards", i)
		}
	})
	t.Run("still-records-when-action-errors", func(t *testing.T) {
		obs, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		myErr := errors.New("boom")
		fn := xwrapcall.New[context.Context]().
			ObservingEvents(obs).
			Do(func(c context.Context) error { return myErr })
		err := fn(context.Background())
		require.ErrorIs(t, err, myErr)
		require.Len(t, view.Events(), 2)
	})
	t.Run("still-records-on-early-return", func(t *testing.T) {
		obs, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		myErr := errors.New("boom")
		earlyReturn := func(c context.Context, next xwrapcall.Next[context.Context]) (context.Context, error) {
			return c, myErr
		}
		fn := xwrapcall.New[context.Context]().
			With(earlyReturn).
			OnEarlyReturnDo(xwrapcall.IgnoreEarlyReturn[context.Context]).
			ObservingEvents(obs).
			Do(func(c context.Context) error { return nil })
		err := fn(context.Background())
		require.ErrorIs(t, err, myErr)

		events := view.Events()
		require.Len(t, events, 2)
		require.Equal(t, xwrapcallobs.FrameEnter, events[0].Kind)
		require.Equal(t, xwrapcallobs.FrameLeave, events[1].Kind)
	})
	t.Run("events-returns-snapshot", func(t *testing.T) {
		obs, view := xwrapcallobs.NewTimings[context.Context](xtime.StdClock)
		fn := xwrapcall.New[context.Context]().
			ObservingEvents(obs).
			Do(func(c context.Context) error { return nil })
		_ = fn(context.Background())

		e1 := view.Events()
		e1[0].StackIndex = 999
		e2 := view.Events()
		require.Equal(t, xwrapcall.StackInd(0), e2[0].StackIndex)
	})
	t.Run("custom-clock", func(t *testing.T) {
		fakeNow := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		fake := &staticClock{now: fakeNow}
		obs, view := xwrapcallobs.NewTimings[context.Context](fake)
		fn := xwrapcall.New[context.Context]().
			ObservingEvents(obs).
			Do(func(c context.Context) error { return nil })
		_ = fn(context.Background())

		events := view.Events()
		require.Len(t, events, 2)
		require.Equal(t, fakeNow, events[0].Ts)
		require.Equal(t, fakeNow, events[1].Ts)
	})
}

type staticClock struct {
	now time.Time
}

func (c *staticClock) Now() time.Time { return c.now }
