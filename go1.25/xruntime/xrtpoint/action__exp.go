package xrtpoint

import (
	"context"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
)

type Context interface {
	context.Context
	ThisPoint() PointConst
	Panic()
	// TODO: Wait(Waitable), Yield()
}

type Action = xwrapcall.SilentAction[Context]
type Middleware = xwrapcall.SilentMiddleware[Context]
type Next = xwrapcall.SilentNext[Context]

var NoopAction Action = func(c Context) {}

func newContext(c context.Context, p PointConst) Context {
	return &ctx{Context: c, p: p}
}

type ctx struct {
	context.Context
	p PointConst
}

func (c *ctx) ThisPoint() PointConst {
	return c.p
}

func (c *ctx) Panic() {
	panic("Injected runtime point caused panic")
}
