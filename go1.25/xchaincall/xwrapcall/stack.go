package xwrapcall

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xshould"
)

// Stack is single threaded.
// Only consecutive Invoke calls are allowed (not parallel).
// Note that Invoke method has signature of Action type.
type Stack[CtxT Context] interface {
	// Invoke executes a call stack.
	// It is NOT safe to call Invoke concurrently.
	Invoke(CtxT) error
	nextAt(StackInd) (next func(CtxT) (CtxT, error), nextCalls *atomic.Uint32)
	// TODO: Reset(), Clone() ?
}

func newStack[CtxT Context](
	mws []Middleware[CtxT],
	a Action[CtxT],
	erA EarlyReturnAction[CtxT],
	obs Observer[CtxT],
) Stack[CtxT] {
	s := &stack[CtxT]{
		mws: mws,
		a:   a,
		erA: erA,
		obs: obs,

		ctxIsCopyable: xshould.TypeImplements[CtxT, CopyableContext[CtxT]]() == nil,
	}
	return s
}

type stack[CtxT Context] struct {
	mws []Middleware[CtxT]
	a   Action[CtxT]
	erA EarlyReturnAction[CtxT]
	obs Observer[CtxT]

	ctxIsCopyable bool
}

var _ Stack[context.Context] = (*stack[context.Context])(nil)

// TODO: support calling Invoke concurrently (solve issue of how distinguish execs from different Invokes, or should not distinguish)
// possible impl:
// - next, _ := tcs.nextAt(0, NextWithNewExecUnit())
// - support hierarchical structure of exec units (as a tree)
func (tcs *stack[CtxT]) Invoke(c CtxT) error {
	if tcs.ctxIsCopyable {
		rt := &stackRuntime{}
		rt.callerStackInd.Store(InvalidStackInd)
		c = tcs.injectRuntime(c, rt)
	}
	next, _ := tcs.nextAt(0)
	_, err := next(c)
	return err
}

func (tcs *stack[CtxT]) nextAt(i StackInd) (func(CtxT) (CtxT, error), *atomic.Uint32) {
	var calls atomic.Uint32
	next := func(c CtxT) (retC CtxT, retErr error) {
		if calls.Add(1) > 1 {
			// TODO: allow multiple next calls, but
			// 1) resolve issue with active stack ind
			// 2) add ability to configure behaviour (e.g. forbid)
			// possible impl for 1):
			// 	 - if calls.Add(1) > 1: execInd := stack.registerNewExecUnit()
			// 	 - next, nextCalls := nextAt(StackInd, NextWithExecInd(execInd)) // default: execInd=0
			//   - stack.ActiveInds() []atomic.Int64
			// possible impl for 2):
			//   - New().OnRepeatedNextCallDo(...)
			//   // should return bool (indicate whether to continue or return)
			//   // unlike PreApply for stage middleware, it is not obvious what will happen
			//   // if return values are (CtxT, error)
			panic("multiple next calls not supported yet")
		}
		if ac, ok := Context(c).(AbortableContext); ok {
			if ac.Aborted() {
				panic(ErrNextAfterAbort)
			}
		}

		var frame *frameEvent[CtxT]
		if tcs.obs != nil {
			frame = &frameEvent[CtxT]{
				execInd:  0,
				stackInd: i,
			}
			tcs.obs.OnFrameEnter(frame)
			defer func() {
				frame.err = retErr
				tcs.obs.OnFrameLeave(frame)
			}()
		}

		var prevRt Runtime
		if tcs.ctxIsCopyable {
			prevRt = ctxCallerRuntime(c)
			nextRt := prevRt.(*stackRuntime).clone()
			nextRt.callerStackInd.Store(i)
			c = tcs.injectRuntime(c, nextRt)
		}

		var err error
		if i < int64(len(tcs.mws)) {
			next, nextCalls := tcs.nextAt(i + 1)
			if frame != nil {
				frame.panicked = true
			}
			c, err = tcs.mws[i](c, next)
			if frame != nil {
				frame.panicked = false
			}
			calledNext := nextCalls.Load() > 0
			if frame != nil {
				frame.calledNext = calledNext
			}
			if !calledNext {
				// TODO: may be issue with force cancel middleware
				// (e.g. when timeout middleware forcefully returns,
				// but goroutine haven't started and therefore next call count
				// haven't increased)
				// possible sol:
				// - public options for next calls (e.g. next(c, NextInParallel()))
				//   - this next call will block until callInd incremented and then start goroutine andreturn
				eai := &earlyReturnInfo[CtxT]{
					stackInd: i,
					mw:       tcs.mws[i],
					rctx:     c,
					rerr:     err,
				}
				c, err = tcs.erA(eai)
			}
		} else if i == int64(len(tcs.mws)) {
			if frame != nil {
				frame.panicked = true
			}
			err = tcs.a(c)
			if frame != nil {
				frame.panicked = false
				frame.calledNext = true
			}
		} else {
			panic(ErrStackOverflow)
		}

		if tcs.ctxIsCopyable {
			c = tcs.injectRuntime(c, prevRt)
		}
		retC, retErr = c, err
		return
	}
	return next, &calls
}

func (tcs *stack[CtxT]) injectRuntime(c CtxT, rt Runtime) CtxT {
	mc := any(c).(CopyableContext[CtxT])
	return mc.CopyOnto(ctxWithCallerRuntime(mc.StdContext(), rt))
}

func getFuncFullname(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

type StackInd = int64

var (
	InvalidStackInd StackInd = -1
)

var ErrNextAfterAbort = errors.New("next was called after Abort()")
var ErrStackOverflow = errors.New("xwrapcall stack overflowed: next was called too many times")
