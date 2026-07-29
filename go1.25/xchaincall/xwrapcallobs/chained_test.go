package xwrapcallobs_test

import (
	"context"
	"testing"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcallobs"
	"github.com/stretchr/testify/require"
)

type recordingObserver struct {
	enters int
	leaves int
}

func (r *recordingObserver) OnFrameEnter(xwrapcall.FrameEnterEvent[context.Context]) {
	r.enters++
}

func (r *recordingObserver) OnFrameLeave(xwrapcall.FrameLeaveEvent[context.Context]) {
	r.leaves++
}

func TestChained(t *testing.T) {
	t.Run("empty/returns-noop", func(t *testing.T) {
		obs := xwrapcallobs.Chained[context.Context]()
		fn := xwrapcall.New[context.Context]().
			ObservingEvents(obs).
			Do(func(c context.Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)
	})
	t.Run("single/returns-it-directly", func(t *testing.T) {
		obs1 := &recordingObserver{}
		obs := xwrapcallobs.Chained[context.Context](obs1)
		require.Same(t, obs1, obs)
	})
	t.Run("multiple/fans-out", func(t *testing.T) {
		obs1 := &recordingObserver{}
		obs2 := &recordingObserver{}
		fn := xwrapcall.New[context.Context]().
			ObservingEvents(xwrapcallobs.Chained(obs1, obs2)).
			Do(func(c context.Context) error { return nil })
		err := fn(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, obs1.enters)
		require.Equal(t, 1, obs1.leaves)
		require.Equal(t, 1, obs2.enters)
		require.Equal(t, 1, obs2.leaves)
	})
}
