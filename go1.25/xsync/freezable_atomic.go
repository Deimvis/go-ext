package xsync

import (
	"sync"

	"github.com/Deimvis/go-ext/go1.25/xatomic"
)

type FreezableAtomic[T any] interface {
	xatomic.Atomic[T]
	Freezable[T]
}

type freezableAtomic[T comparable, A xatomic.Atomic[T]] struct {
	v        A
	freezeCb func()
	mu       sync.RWMutex
}

var _ FreezableAtomic[any] = &freezableAtomic[any, xatomic.Atomic[any]]{}

func (fa *freezableAtomic[T, A]) Load() T {
	return fa.v.Load()
}

func (fa *freezableAtomic[T, A]) Store(v T) {
	fa.mu.RLock()
	defer fa.mu.RUnlock()
	if fa.freezeCb != nil {
		fa.freezeCb()
	}
	fa.v.Store(v)
}

func (fa *freezableAtomic[T, A]) Swap(new T) (old T) {
	fa.mu.RLock()
	defer fa.mu.RUnlock()
	if fa.freezeCb != nil {
		fa.freezeCb()
	}
	return fa.v.Swap(new)
}

func (fa *freezableAtomic[T, A]) CompareAndSwap(old T, new T) bool {
	fa.mu.RLock()
	defer fa.mu.RUnlock()
	if fa.freezeCb != nil {
		fa.freezeCb()
	}
	return fa.v.CompareAndSwap(old, new)
}

func (fa *freezableAtomic[T, A]) Freeze(violationCallback func()) bool {
	fa.mu.Lock()
	defer fa.mu.Unlock()
	if fa.freezeCb != nil {
		return false
	}
	fa.freezeCb = violationCallback
	return true
}

func (fa *freezableAtomic[T, A]) Freezed() bool {
	fa.mu.RLock()
	defer fa.mu.RUnlock()
	return fa.freezeCb != nil
}
