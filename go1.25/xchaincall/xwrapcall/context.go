package xwrapcall

import "github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcallctx"

type Context = xwrapcallctx.Context

// CopyableContext adds ability to
// access context.Context directly
// and rebase a copy onto new context.Context.
//
// CopyableContext is required for
// injecting Runtime info.
type CopyableContext[Self any] = xwrapcallctx.WithCopying[Self]

// AbortableContext adds new semantic state for call.
// With Context the only state available is either
// error or nothing happened.
// Abort helps to indicate that call was interrupted
// before completing work it was supposed to do.
// Calling next after Abort() is called causes panic.
//
// Abort is not necessarily error and it's a valid case
// when call was aborted with nil error returned.
// And non-nil error returned does not necessarily means abort,
// it may indicate only partial work was done.
//
// If call is not supposed to abort and any abort
// is an error, then AbortContext is not needed
// and early return with an error may be used.
//
// The common rule about how to use Abort, error returned and panic:
// - Work done: no error, no abort
// - Expected error occured and work is not done (e.g. timeout): error, abort
// - Expected error occured while doing work (e.g. work is only partially completed): error, no abort
// - Expected reason why work should not be done (e.g. idempotent call and work is alredy done): no error, abort
// - Unexpected, uncontrollable error occured: panic
type AbortableContext = xwrapcallctx.WithAbort
