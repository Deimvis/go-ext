package xsync

import (
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xatomic"
)

func NewFreezableInt32() FreezableAtomicInt[int32] {
	a := &freezableAtomicInt[int32, xatomic.AtomicInt[int32]]{}
	a.v = &atomic.Int32{}
	return a
}

func NewFreezableInt64() FreezableAtomicInt[int64] {
	a := &freezableAtomicInt[int64, xatomic.AtomicInt[int64]]{}
	a.v = &atomic.Int64{}
	return a
}

func NewFreezableUint64() FreezableAtomicInt[uint64] {
	a := &freezableAtomicInt[uint64, xatomic.AtomicInt[uint64]]{}
	a.v = &atomic.Uint64{}
	return a
}

type FreezableAtomicInt[T int32 | int64 | uint32 | uint64 | uintptr] interface {
	xatomic.AtomicInt[T]
	Freezable[T]
}

type freezableAtomicInt[T int32 | int64 | uint32 | uint64 | uintptr, A xatomic.AtomicInt[T]] struct {
	freezableAtomic[T, A]
}

var _ FreezableAtomicInt[int32] = &freezableAtomicInt[int32, xatomic.AtomicInt[int32]]{}

func (fai *freezableAtomicInt[T, A]) Add(delta T) (new T) {
	if fai.freezeCb != nil {
		fai.freezeCb()
	}
	return fai.v.Add(delta)
}

func (fai *freezableAtomicInt[T, A]) And(mask T) (old T) {
	if fai.freezeCb != nil {
		fai.freezeCb()
	}
	return fai.v.And(mask)
}

func (fai *freezableAtomicInt[T, A]) Or(mask T) (old T) {
	if fai.freezeCb != nil {
		fai.freezeCb()
	}
	return fai.v.Or(mask)
}
