//go:build playground

package playground

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

// TODO: sampler must be an interface (not only sampling by min interval)
func WithSampling(fn func(), minInterval time.Duration) (func(), error) {
	var samplingFn func()

	minIntervalS := int64(minInterval.Seconds())
	xmust.NotEq(minIntervalS, 0, "min interval between logs must be > 0")
	if minIntervalS == 0 {
		return samplingFn, errors.New("min interval between logs must be > 0")
	}

	var prevLogTs atomic.Int64
	prevLogTs.Store(-int64(minIntervalS))

	samplingFn = func(fn func(), prevLogTs *atomic.Int64) func() {
		return func() {
			now := nowFn().Unix()
			prevLogTsSnapshot := prevLogTs.Load()
			for now-prevLogTsSnapshot >= minIntervalS {
				if prevLogTs.CompareAndSwap(prevLogTsSnapshot, now) {
					fn()
					break
				} else {
					prevLogTsSnapshot = prevLogTs.Load()
				}
			}
		}
	}(fn, &prevLogTs)
	return samplingFn, nil
}

var nowFn = time.Now
