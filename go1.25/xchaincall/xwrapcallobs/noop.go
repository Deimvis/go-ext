package xwrapcallobs

import "github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"

func Noop[CtxT xwrapcall.Context]() xwrapcall.Observer[CtxT] {
	return noop[CtxT]{}
}

type noop[CtxT xwrapcall.Context] struct{}

func (noop[CtxT]) OnFrameEnter(xwrapcall.FrameEnterEvent[CtxT]) {}
func (noop[CtxT]) OnFrameLeave(xwrapcall.FrameLeaveEvent[CtxT]) {}
