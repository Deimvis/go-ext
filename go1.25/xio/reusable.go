package xio

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

// ReusableReadSeekerFactory is helpful for reusing io.ReadSeeker,
// by ensuring that no one will access it concurrently.
// Note that calling New() will return new io.ReadSeekCloser and Close previous one.
//
// Common issue with raw io.ReadSeeker is that when attempting to reuse
// and moving cursor to the start, concurrent access is io.ReadSeeker is possible
// which leads to race conditions and errors.
// For example, http.Request.Body has no option to track its usage and
// any concurrent Read() calls are possible anytime later,
// which should return error if "http.Body.Request is no longer available",
// and in order to indicate that, http.Body.Request should be closed
// (after waiting all current read requests to finish) and only then io.ReadSeeker may be reused,
// but keeping http.Body.Request return error on Read() from now on.
type ReusableReadSeekerFactory interface {
	// New calls Close() on previously returned io.ReadSeekCloser.
	// Close() call ensures that underlying io.ReadSeeker is no longer used.
	// No more than one instance of non-closed io.ReadSeekCloser is available at a time.
	New() io.ReadSeekCloser
}

// implementation draft:

type reusableReadSeekerFactory struct {
	prev atomic.Pointer[concurrentReadSeekCloser]
	// TODO: resolve efficiently (who gets non-closed instance) when concurrent calls to New() ?
	//       upd: maybe better forbid concurrent calls to New(), because this case is itself invalid.
}

// concurrentReadSeekCloser allows to Close() concurrently
// (Read() and Seek() are not synchronized and MUST NOT be called concurrently).
// Safe to Close() multiple times and concurrently.
// Read() and Seek() return errClosed even if closing is in progress.
//
// Implementation notes:
//   - implementation with rwmutex would block Read() during Close(),
//     with current one we can indicate closing and exit early
//   - implementation with rwmutex would use 2 operations (inc/desc) for Read() on closed,
//     but we can leverage that once closed, all Read() calls will return errClosed.
type concurrentReadSeekCloser struct {
	rs io.ReadSeeker

	closed atomic.Int32 // 0: not closed, 1: closing, 2: closed
	// users include readers and seekers
	users      atomic.Int32
	usersWait  atomic.Int32
	usersDone  sync.Cond
	closerDone sync.Cond
}

func (rsc *concurrentReadSeekCloser) Read(p []byte) (n int, err error) {
	rsc.use(func() {
		n, err = rsc.rs.Read(p)
	})
	return
}

func (rsc *concurrentReadSeekCloser) Seek(offset int64, whence int) (newOffset int64, err error) {
	rsc.use(func() {
		newOffset, err = rsc.rs.Seek(offset, whence)
	})
	return
}

func (rsc *concurrentReadSeekCloser) Close() error {
	old := rsc.closed.Load()
	for old == 0 {
		if rsc.closed.CompareAndSwap(old, 1) {
			break
		}
		old = rsc.closed.Load()
	}
	if old == 2 {
		// already closed
		return nil
	}
	if old == 1 {
		// closing
		func() {
			rsc.closerDone.L.Lock()
			defer rsc.closerDone.L.Unlock()

			for rsc.closed.Load() != 2 {
				rsc.closerDone.Wait()
			}
		}()
		return nil
	}

	r := rsc.users.Add(-maxRscUsers) + maxRscUsers
	if r != 0 && rsc.usersWait.Add(r) != 0 {
		func() {
			rsc.usersDone.L.Lock()
			defer rsc.usersDone.L.Unlock()

			for rsc.usersWait.Load() != 0 {
				rsc.usersDone.Wait()
			}
		}()
	}

	rsc.users.Add(maxRscUsers)
	func() {
		// wake up concurrent Close()
		rsc.closerDone.L.Lock()
		defer rsc.closerDone.L.Unlock()

		rsc.closerDone.Broadcast()
	}()

	return nil
}

func (rsc *concurrentReadSeekCloser) use(fn func()) error {
	if rsc.closed.Load() > 0 {
		// closed/closing
		return errClosed
	}
	if rsc.users.Add(1) < 0 {
		// closed/closing: just started closing
		rsc.users.Add(-1)
		return errClosed
	}

	fn()

	if rsc.users.Add(-1) < 0 {
		func() {
			rsc.usersDone.L.Lock()
			defer rsc.usersDone.L.Unlock()
			if rsc.usersWait.Add(-1) == 0 {
				rsc.usersDone.Signal()
			}
		}()
	}

	return nil
}

const (
	maxRscUsers = 1 << 30
)

var (
	errClosed = errors.New("closed")
)
