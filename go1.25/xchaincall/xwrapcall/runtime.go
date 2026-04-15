package xwrapcall

import (
	"context"
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

type Runtime interface {
	CallerStackIndex() StackInd
}

// CallerRuntime returns caller's stack runtime.
func CallerRuntime(ctx context.Context) (Runtime, bool) {
	v := ctx.Value(callerRuntimeCtxKey{})
	if v == nil {
		return nil, false
	}
	return v.(Runtime), true
}

func MustCallerRuntime(ctx context.Context) Runtime {
	return xmust.Ok(CallerRuntime(ctx))
}

type stackRuntime struct {
	callerStackInd atomic.Int64
}

var _ Runtime = (*stackRuntime)(nil)

func (sr *stackRuntime) CallerStackIndex() StackInd {
	return sr.callerStackInd.Load()
}

func (sr *stackRuntime) clone() *stackRuntime {
	srCopy := &stackRuntime{}
	srCopy.callerStackInd.Store(sr.callerStackInd.Load())
	return srCopy
}

func ctxCallerRuntime(ctx context.Context) Runtime {
	return ctx.Value(callerRuntimeCtxKey{}).(Runtime)
}

func ctxWithCallerRuntime(ctx context.Context, rt Runtime) context.Context {
	return context.WithValue(ctx, callerRuntimeCtxKey{}, rt)
}

type callerRuntimeCtxKey struct{}
