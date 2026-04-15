package xwrapcallctx

import (
	"context"
)

func New(ctx context.Context) *impl {
	return &impl{Context: ctx}
}

// impl is thread-safe implementation of xwrapcall.AbortableContext.
// impl forbids Abort() from being called more than once.
type impl struct {
	context.Context
	HavingAbort
}
