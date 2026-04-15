package xio

import (
	"sync/atomic"
)

type CountingReaderWrapState interface {
	BytesRead() int64
}

// CountingReaderWrap calculates total bytes read on ReadFn calls.
// Wrap is thread-safe.
func CountingReaderWrap(fn ReadFn) (ReadFn, CountingReaderWrapState) {
	state := &countingReaderWrapState{}
	wrappedFn := func(p []byte) (int, error) {
		n, err := fn(p)
		if n > 0 {
			state.bytesRead.Add(int64(n))
		}
		return n, err
	}
	return wrappedFn, state
}

type countingReaderWrapState struct {
	bytesRead atomic.Int64
}

func (c *countingReaderWrapState) BytesRead() int64 {
	return c.bytesRead.Load()
}
