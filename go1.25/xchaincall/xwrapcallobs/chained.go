package xwrapcallobs

import "github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"

func Chained[CtxT xwrapcall.Context](observers ...xwrapcall.Observer[CtxT]) xwrapcall.Observer[CtxT] {
	if len(observers) == 0 {
		return Noop[CtxT]()
	}
	if len(observers) == 1 {
		return observers[0]
	}
	return chained[CtxT](observers)
}

type chained[CtxT xwrapcall.Context] []xwrapcall.Observer[CtxT]

func (c chained[CtxT]) OnFrameEnter(info xwrapcall.FrameEnterEvent[CtxT]) {
	for _, o := range c {
		o.OnFrameEnter(info)
	}
}

func (c chained[CtxT]) OnFrameLeave(info xwrapcall.FrameLeaveEvent[CtxT]) {
	for _, o := range c {
		o.OnFrameLeave(info)
	}
}
