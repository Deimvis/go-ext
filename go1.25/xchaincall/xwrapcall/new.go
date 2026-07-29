package xwrapcall

import (
	"context"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

func New[CtxT Context]() Builder[CtxT] {
	return &builder[CtxT]{}
}

type Builder[CtxT Context] interface {
	With(...Middleware[CtxT]) Builder[CtxT]
	// OnEarlyReturnDo allows to customize action called
	// when middleware returned before calling next.
	OnEarlyReturnDo(EarlyReturnAction[CtxT]) Builder[CtxT]
	Do(Action[CtxT]) Action[CtxT]

	// ObservingEvents allows to observe events synchronously
	// with call stack execution.
	ObservingEvents(Observer[CtxT]) Builder[CtxT]

	// TODO: allow dynamically customize stack
	// (OnStack(func(...) Stack))
	// (maybe someone would like to link middlewares directly and sacrifice transparency for performance)
}

type builder[CtxT Context] struct {
	mws []Middleware[CtxT]
	a   Action[CtxT]
	erA EarlyReturnAction[CtxT]
	obs Observer[CtxT]
}

var _ Builder[context.Context] = (*builder[context.Context])(nil)

func (b *builder[CtxT]) With(mws ...Middleware[CtxT]) Builder[CtxT] {
	b.mws = mws
	return b
}

func (b *builder[CtxT]) OnEarlyReturnDo(erA EarlyReturnAction[CtxT]) Builder[CtxT] {
	b.erA = erA
	return b
}

func (b *builder[CtxT]) ObservingEvents(o Observer[CtxT]) Builder[CtxT] {
	b.obs = o
	return b
}

func (b *builder[CtxT]) Do(a Action[CtxT]) Action[CtxT] {
	b.a = a
	return b.build()
}

func (b *builder[CtxT]) build() Action[CtxT] {
	xmust.NotNilInterface(b.a, "action not set")
	if b.erA == nil {
		b.erA = earlyReturnAction_default[CtxT]
	}
	s := newStack(b.mws, b.a, b.erA, b.obs)
	return s.Invoke
}

type earlyReturnInfo[CtxT Context] struct {
	stackInd StackInd
	mw       Middleware[CtxT]
	rctx     CtxT
	rerr     error
}

var _ EarlyReturnInfo[Context] = &earlyReturnInfo[Context]{}

func (eai *earlyReturnInfo[CtxT]) StackIndex() StackInd {
	return eai.stackInd
}

func (eai *earlyReturnInfo[CtxT]) Middleware() Middleware[CtxT] {
	return eai.mw
}

func (eai *earlyReturnInfo[CtxT]) ReturnedValues() (CtxT, error) {
	return eai.rctx, eai.rerr
}
