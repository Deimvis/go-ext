package xtesting

import (
	"fmt"
	"reflect"
	"sync"
)

func NewBarrier() (BarrierParker, BarrierGuard) {
	b := &barrier{}
	b.parkCond = sync.NewCond(&b.parkMu)
	b.waiters = nil
	b.disabled = false
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

	// Disable flushes parked units
	// and returns their count.
	// Any subsequent Park calls
	// will be no-op until barrier is enabled again.
	// And subsequent Wait, Pass calls
	// will panic.
	Disable() int
	// Enable enables barrier.
	// If barrier is not disabled, it is no-op.
	Enable()
}

type barrierGuardSugar interface {
	WaitOne()
	WaitN(count int)

	Count(pred func(tags []any) bool) int
	// CountHavingTags counts only ones
	// that include all given tags.
	CountHavingTags(tags ...any) int
	CountAll() int

	PassOne()
	PassN(count int)

	// NOTE: maybe add Flush <=> PassN(CountAll())
}

// BarrierGuard implements wait-inspect-pass loop.
type BarrierGuard interface {
	BarrierGuardBasis
	barrierGuardSugar
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

	enqueued := func() bool {
		bp.b.parkMu.Lock()
		defer bp.b.parkMu.Unlock()
		if bp.b.disabled {
			return false
		}
		bp.b.waiters = append(bp.b.waiters, w)
		bp.b.parkCond.Broadcast()
		return true
	}()
	if !enqueued {
		return
	}

	<-w.release
}

type barrierGuard struct {
	b *barrier
}

var _ BarrierGuard = (*barrierGuard)(nil)

func (bg *barrierGuard) Wait(pred func(tags []any) bool, count int) {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	for !bg.b.disabled && bg.count_locked(pred) < count {
		bg.b.parkCond.Wait()
	}
	if bg.b.disabled {
		panic("barrier is disabled")
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
	for !bg.b.disabled && bg.count_locked(pred) < count {
		bg.b.parkCond.Wait()
	}
	if bg.b.disabled {
		panic("barrier is disabled")
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

func (bg *barrierGuard) Disable() int {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	n := bg.flush_locked(func([]any) bool { return true })
	bg.b.disabled = true
	bg.b.parkCond.Broadcast()
	return n
}

func (bg *barrierGuard) Enable() {
	bg.b.parkMu.Lock()
	defer bg.b.parkMu.Unlock()
	bg.b.disabled = false
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
		return fmt.Errorf("park tag %v of type %T is not comparable", t, t)
	}
	return nil
}

type barrier struct {
	parkMu   sync.Mutex
	parkCond *sync.Cond
	waiters  []*waiter // FIFO order of arrival
	disabled bool
}

type waiter struct {
	tags    []any
	release chan struct{} // closed by the guard to unpark
}
