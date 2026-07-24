package xtesting

import (
	"fmt"
	"reflect"
	"sync"
)

// TODO: debug, add Gate (knows who passed through, but always open)

func NewBarrier() (BarrierParker, BarrierGuard) {
	b := &barrier{}
	b.parkCond = sync.NewCond(&b.parkMu)
	b.waiters = nil
	b.closed = false
	return &barrierParker{b}, &barrierGuard{b}
}

type BarrierParker interface {
	Park(tags ...any)
}

type BarrierGuardBasis interface {
	Wait(pred func(tags []any) bool, count int)
	Range(f func(tags []any) bool)
	// Pass blocks until at least given count
	// of parked units match pred.
	Pass(pred func(tags []any) bool, count int)

	// NOTE: could not come up with
	// reasonable examples so far
	// why behaviour below
	// may be really needed.
	// -
	// // PassAll immediately releases everyone who has already parked.
	// PassAll() int
	// // sugared
	// PassFn(pred) int
	// UPD: better naming: Flush & FlushHavingTags

	// Close flushes parked units
	// and any future or not-finished
	// Park, Wait, Pass calls
	// will panic.
	Close() int
}

// BarrierGuard implements wait-inspect-pass loop.
type BarrierGuard interface {
	BarrierGuardBasis

	// -- sugar (shortcuts) --

	WaitOne()
	WaitN(count int)

	Count(pred func(tags []any) bool) int
	// CountHavingTags counts only ones
	// that include all given tags.
	CountHavingTags(tags ...any) int
	CountAll() int

	PassOne()
	PassN(count int)
}

type barrierParker struct {
	b *barrier
}

var _ BarrierParker = (*barrierParker)(nil)

func (bp *barrierParker) Park(tags ...any) {
	for _, t := range tags {
		err := validateTag(t)
		if err != nil {
			panic(err)
		}
	}
	w := &waiter{tags: tags, release: make(chan struct{})}

	func() {
		bp.b.parkMu.Lock()
		defer bp.b.parkMu.Unlock()
		if bp.b.closed {
			panic("xtesting: barrier is closed")
		}
		bp.b.waiters = append(bp.b.waiters, w)
		bp.b.parkCond.Broadcast()
	}()

	<-w.release
}

type barrierGuard struct {
	b *barrier
}

var _ BarrierGuard = (*barrierGuard)(nil)

func (bg *barrierGuard) Wait(pred func(tags []any) bool, count int) {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	for !bg.b.closed && bg.count_locked(pred) < count {
		bg.b.parkCond.Wait()
	}
	if bg.b.closed {
		panic("xtesting: barrier is closed")
	}
}

func (bg *barrierGuard) WaitOne() {
	bg.Wait(func([]any) bool { return true }, 1)
}

func (bg *barrierGuard) WaitN(count int) {
	bg.Wait(func([]any) bool { return true }, count)
}

func (bg *barrierGuard) Range(f func(tags []any) bool) {
	bg.b.parkMu.Lock()
	snapshot := make([][]any, len(bg.b.waiters))
	for i, w := range bg.b.waiters {
		snapshot[i] = append([]any(nil), w.tags...)
	}
	bg.b.parkMu.Unlock()

	for _, tags := range snapshot {
		if !f(tags) {
			return
		}
	}
}

func (bg *barrierGuard) Count(pred func(tags []any) bool) int {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	return bg.count_locked(pred)
}

func (bg *barrierGuard) CountHavingTags(tags ...any) int {
	pred := func(actTags []any) bool {
		for _, exp := range tags {
			found := false
			for _, act := range actTags {
				if act == exp {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	return bg.Count(pred)
}

func (bg *barrierGuard) CountAll() int {
	return bg.Count(func([]any) bool { return true })
}

func (bg *barrierGuard) Pass(pred func(tags []any) bool, count int) {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	for !bg.b.closed && bg.count_locked(pred) < count {
		bg.b.parkCond.Wait()
	}
	if bg.b.closed {
		panic("xtesting: barrier is closed")
	}
	released := 0
	keepNextInd := 0
	for _, w := range bg.b.waiters {
		if released < count && pred(w.tags) {
			close(w.release)
			released++
		} else {
			bg.b.waiters[keepNextInd] = w
			keepNextInd++
		}
	}
	clear(bg.b.waiters[keepNextInd:]) // gc
	bg.b.waiters = bg.b.waiters[:keepNextInd]
}

func (bg *barrierGuard) PassOne() {
	bg.Pass(func([]any) bool { return true }, 1)
}

func (bg *barrierGuard) PassN(count int) {
	bg.Pass(func([]any) bool { return true }, count)
}

func (bg *barrierGuard) Close() int {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	n := bg.flush_locked(func([]any) bool { return true })
	bg.b.closed = true
	bg.b.parkCond.Broadcast()
	return n
}

func (bg *barrierGuard) count_locked(pred func(tags []any) bool) int {
	n := 0
	for _, w := range bg.b.waiters {
		if pred(w.tags) {
			n++
		}
	}
	return n
}

func (bg *barrierGuard) flush_locked(pred func(tags []any) bool) int {
	n := 0
	keepNextInd := 0
	for _, w := range bg.b.waiters {
		if pred(w.tags) {
			close(w.release)
			n++
		} else {
			bg.b.waiters[keepNextInd] = w
			keepNextInd++
		}
	}
	clear(bg.b.waiters[keepNextInd:]) // gc
	bg.b.waiters = bg.b.waiters[:keepNextInd]
	return n
}

func validateTag(t any) error {
	if t == nil || !reflect.TypeOf(t).Comparable() {
		return fmt.Errorf("xtesting: park tag %v of type %T is not comparable", t, t)
	}
	return nil
}

type barrier struct {
	parkMu   sync.Mutex
	parkCond *sync.Cond
	waiters  []*waiter // FIFO order of arrival
	closed   bool
}

type waiter struct {
	tags    []any
	release chan struct{} // closed by the guard to unpark
}
