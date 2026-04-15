package xrtpointtest

import "testing"

// TODO: move to another package as generic approach for test callback

type TestCallback func(t *testing.T)

func RunCallbacks(t *testing.T, cbs ...TestCallback) {
	if t.Failed() {
		return
	}
	if r := recover(); r != nil {
		panic(r)
	}
	for _, cb := range cbs {
		cb(t)
	}
}
