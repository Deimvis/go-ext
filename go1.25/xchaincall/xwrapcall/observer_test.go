package xwrapcall

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type recordingObserver[CtxT Context] struct {
	events []observerEvent[CtxT]
}

type observerEvent[CtxT Context] struct {
	kind       string
	stackInd   StackInd
	err        error
	calledNext bool
}

func (r *recordingObserver[CtxT]) OnFrameEnter(info FrameEnterEvent[CtxT]) {
	r.events = append(r.events, observerEvent[CtxT]{
		kind:     "enter",
		stackInd: info.StackIndex(),
	})
}

func (r *recordingObserver[CtxT]) OnFrameLeave(info FrameLeaveEvent[CtxT]) {
	r.events = append(r.events, observerEvent[CtxT]{
		kind:       "leave",
		stackInd:   info.StackIndex(),
		err:        info.Error(),
		calledNext: info.CalledNext(),
	})
}

func TestObserver(t *testing.T) {
	t.Run("no-mws/action-only", func(t *testing.T) {
		obs := &recordingObserver[Context]{}
		fn := New[Context]().
			ObservingEvents(obs).
			Do(func(c Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)
		require.Equal(t, []observerEvent[Context]{
			{kind: "enter", stackInd: 0},
			{kind: "leave", stackInd: 0, calledNext: true},
		}, obs.events)
	})
	t.Run("mws/order/enter-descending_leave-ascending", func(t *testing.T) {
		obs := &recordingObserver[Context]{}
		passThrough := func(c Context, next Next[Context]) (Context, error) {
			return next(c)
		}
		fn := New[Context]().
			With(passThrough, passThrough, passThrough).
			ObservingEvents(obs).
			Do(func(c Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)
		expKinds := []string{"enter", "enter", "enter", "enter", "leave", "leave", "leave", "leave"}
		expInds := []StackInd{0, 1, 2, 3, 3, 2, 1, 0}
		actKinds := make([]string, len(obs.events))
		actInds := make([]StackInd, len(obs.events))
		for i, e := range obs.events {
			actKinds[i] = e.kind
			actInds[i] = e.stackInd
		}
		require.Equal(t, expKinds, actKinds)
		require.Equal(t, expInds, actInds)
	})
	t.Run("mws/error-propagated-in-leave", func(t *testing.T) {
		obs := &recordingObserver[Context]{}
		myErr := errors.New("boom")
		passThrough := func(c Context, next Next[Context]) (Context, error) {
			return next(c)
		}
		fn := New[Context]().
			With(passThrough).
			ObservingEvents(obs).
			Do(func(c Context) error { return myErr })
		err := fn(context.Background())
		require.ErrorIs(t, err, myErr)
		leaves := []observerEvent[Context]{}
		for _, e := range obs.events {
			if e.kind == "leave" {
				leaves = append(leaves, e)
			}
		}
		require.Len(t, leaves, 2)
		require.ErrorIs(t, leaves[0].err, myErr)
		require.ErrorIs(t, leaves[1].err, myErr)
	})
	t.Run("mws/early-return_calledNext-false", func(t *testing.T) {
		obs := &recordingObserver[Context]{}
		myErr := errors.New("boom")
		earlyReturn := func(c Context, next Next[Context]) (Context, error) {
			return c, myErr
		}
		fn := New[Context]().
			With(earlyReturn).
			OnEarlyReturnDo(IgnoreEarlyReturn[Context]).
			ObservingEvents(obs).
			Do(func(c Context) error { return nil })
		err := fn(context.Background())
		require.ErrorIs(t, err, myErr)
		require.Len(t, obs.events, 2)
		require.Equal(t, "enter", obs.events[0].kind)
		require.Equal(t, StackInd(0), obs.events[0].stackInd)
		require.Equal(t, "leave", obs.events[1].kind)
		require.Equal(t, StackInd(0), obs.events[1].stackInd)
		require.False(t, obs.events[1].calledNext)
	})
	t.Run("mws/panic-still-fires-leave", func(t *testing.T) {
		obs := &recordingObserver[Context]{}
		panicking := func(c Context, next Next[Context]) (Context, error) {
			panic("boom")
		}
		fn := New[Context]().
			With(panicking).
			ObservingEvents(obs).
			Do(func(c Context) error { return nil })
		require.Panics(t, func() {
			_ = fn(context.Background())
		})
		require.Len(t, obs.events, 2)
		require.Equal(t, "enter", obs.events[0].kind)
		require.Equal(t, "leave", obs.events[1].kind)
	})
	t.Run("no-observer/no-panic", func(t *testing.T) {
		fn := New[Context]().
			Do(func(c Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)
	})
	t.Run("abort-context/mw-aborts-without-next/leave-still-fires", func(t *testing.T) {
		obs := &recordingObserver[AbortableContext]{}
		abortMw := func(c AbortableContext, next Next[AbortableContext]) (AbortableContext, error) {
			c.Abort()
			return c, nil
		}
		fn := func(c AbortableContext) error { return nil }
		fn = New[AbortableContext]().
			With(abortMw).
			ObservingEvents(obs).
			Do(fn)
		err := fn(actxNew())
		require.NoError(t, err)
		require.Len(t, obs.events, 2)
		require.Equal(t, "enter", obs.events[0].kind)
		require.Equal(t, "leave", obs.events[1].kind)
		require.False(t, obs.events[1].calledNext)
	})
}
