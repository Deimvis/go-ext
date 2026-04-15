package xrtpointss

import (
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xruntime/xrtpoint"
)

func DoPanicN(n uint64) xrtpoint.Action {
	return DoPanic(func(c xrtpoint.Context, callInd uint64) bool {
		return callInd < n
	})
}

// DoPanicBefore does not provide any synchronization guarantees.
// It's NOT guaranteed that operations are processed sequentially by callInd.
func DoPanicBefore(pred func(c xrtpoint.Context, callInd uint64) bool) xrtpoint.Action {
	var reached atomic.Uint64
	return DoPanic(func(c xrtpoint.Context, callInd uint64) bool {
		if reached.Load() > 0 {
			return false
		}
		if pred(c, callInd) {
			for reached.Load() == 0 {
				reached.CompareAndSwap(0, 1)
			}
			return false
		}
		return true
	})
}

func DoPanic(pred func(c xrtpoint.Context, callInd uint64) bool) xrtpoint.Action {
	var callInd atomic.Uint64
	return func(c xrtpoint.Context) {
		nextInd := callInd.Add(1)
		if pred(c, nextInd-1) {
			c.Panic()
		}
	}
}
