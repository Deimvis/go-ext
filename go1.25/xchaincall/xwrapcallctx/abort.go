package xwrapcallctx

import (
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

type abort interface {
	// TODO: support multiple Aborts? Concurrent aborts? Enriching AbortInfo?
	Abort(...AbortOption)
	Aborted() bool
	AbortInfo() AbortInfo
}

type WithAbort interface {
	Context
	abort
}

// TODO: implement BaseContext which users can embed in their custom context impl
// (maybe xwrapcallctx.Base)

type AbortInfoMutable interface {
	AbortInfo
	SetReason(string)
	SetFields(...Field)
}

type AbortInfo interface {
	Reason() string
	Fields() []Field
}

type Field struct {
	Key   string
	Value any
}

// ContextAbort is thread-safe implementation of xwrapcall.WithAbort.
// ContextAbort forbids Abort() from being called more than once.
type ContextAbort struct {
	aborted   atomic.Int64
	abortInfo atomic.Pointer[abortInfo]
}

var _ abort = (*ContextAbort)(nil)

func (ha *ContextAbort) Abort(opts ...AbortOption) {
	xmust.Eq(ha.aborted.Add(1), 1, "Abort was already called")

	ai := &abortInfo{}
	for _, opt := range opts {
		opt(ai)
	}
	ha.abortInfo.Store(ai)
}

func (ha *ContextAbort) Aborted() bool {
	return ha.aborted.Load() >= 1
}

func (ha *ContextAbort) AbortInfo() AbortInfo {
	return ha.abortInfo.Load()
}
