package xrtpointtest

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xruntime/xrtpoint"
)

func RequireUnreachable(t *testing.T) xrtpoint.Middleware {
	return func(c xrtpoint.Context, next xrtpoint.Next) xrtpoint.Context {
		// TODO: add tracing info?
		t.Log("Unreachable point has been reached")
		t.FailNow()
		return next(c)
	}
}

// TODO: RequireCalled (called at least once) (TODO: naming)

func RequireCalledN(t *testing.T, n uint64, cbs *[]TestCallback) xrtpoint.Middleware {
	xmust.NotNilPtr(cbs, "nil pointer to test callbacks")
	cnt := &atomic.Uint64{}
	cbCalled := &atomic.Uint32{}
	*cbs = append(*cbs, func(t *testing.T) {
		cbCalled.Store(1)
		require.Equal(t, n, cnt.Load())
	})

	return func(c xrtpoint.Context, next xrtpoint.Next) xrtpoint.Context {
		if cbCalled.Load() != 0 {
			t.Log("xrtpoint.Point was called after test callback")
			t.FailNow()
		}
		cnt.Add(1)
		return next(c)
	}
}

// NOTE: to implement this check it requires to use line number or something like this,
// but intentionally I don't want to provide this info from Point. Do not see workaround for this check right now.
// Probably CallLocationId() string in PointConst will help.
//
// func RequireExactlyOnePointMatched(t *testing.T) xrtpoint.Middleware {
// 	var pointId atomic.Pointer[string]
// 	return func(c xrtpoint.Context, next xrtpoint.Next) xrtpoint.Context {
// 		curPointId := c.ThisPoint().CallPackagePath()
// 		for pointId.Load() == nil {
// 			pointId.CompareAndSwap(nil, &curPointId)
// 		}
// 		expPointId := *pointId.Load()
// 		require.Equal(t, expPointId, curPointId)
// 		return next(c)
// 	}
// }
