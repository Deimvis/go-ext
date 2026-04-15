package xwrapcallctx

import (
	"context"
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

type WithAbort interface {
	Abort(...xwrapcall.AbortOption)
	Aborted() bool
	AbortInfo() xwrapcall.AbortInfo
}

func NewHavingAbort() *HavingAbort {
	return &HavingAbort{}
}

// HavingAbort is thread-safe implementation of xwrapcall.WithAbort.
// HavingAbort forbids Abort() from being called more than once.
type HavingAbort struct {
	aborted   atomic.Int64
	abortInfo atomic.Pointer[abortInfo]
}

var _ WithAbort = (*HavingAbort)(nil)

func (ha *HavingAbort) Abort(opts ...xwrapcall.AbortOption) {
	xmust.Eq(ha.aborted.Add(1), 1, "Abort was already called")

	ai := &abortInfo{}
	for _, opt := range opts {
		opt(ai)
	}
	ha.abortInfo.Store(ai)
}

func (ha *HavingAbort) Aborted() bool {
	return ha.aborted.Load() >= 1
}

func (ha *HavingAbort) AbortInfo() xwrapcall.AbortInfo {
	return ha.abortInfo.Load()
}

// interface guard
type _ig_ContextWithAbort struct {
	context.Context
	WithAbort
}

var _ xwrapcall.Context = (*_ig_ContextWithAbort)(nil)
var _ xwrapcall.AbortableContext = (*_ig_ContextWithAbort)(nil)
