package xwrapcallctx

import "context"

// TODO: implement BaseContext which users can embed in their custom context impl
// (maybe xwrapcallctx.Base)

// Context is a minimal context interface
// needed for xwrapcall.
// NOTE: leave as alias in order to
// make implementations (e.g. CopyOnto(Context))
// be compatible with potential third-party interfaces
type Context = context.Context

func New(ctx context.Context) *impl {
	return &impl{Context: ctx}
}

// impl is thread-safe implementation of xwrapcall.AbortableContext.
// impl forbids Abort() from being called more than once.
type impl struct {
	context.Context
	ContextAbort
}
