package xatomic

import "sync/atomic"

type Atomic[T any] interface {
	Load() T
	Store(T)
	Swap(new T) (old T)
	CompareAndSwap(old T, new T) bool
}

var _ Atomic[*any] = &atomic.Pointer[any]{}

type AtomicInt[T int32 | int64 | uint32 | uint64 | uintptr] interface {
	Atomic[T]
	Add(T) (new T)
	And(mask T) (old T)
	Or(mask T) (old T)
}

var _ AtomicInt[int32] = &atomic.Int32{}
var _ AtomicInt[int64] = &atomic.Int64{}
var _ AtomicInt[uint32] = &atomic.Uint32{}
var _ AtomicInt[uint64] = &atomic.Uint64{}
var _ AtomicInt[uintptr] = &atomic.Uintptr{}
