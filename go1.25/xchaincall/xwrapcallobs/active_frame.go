package xwrapcallobs

import (
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
)

func NewActiveFrame[CtxT xwrapcall.Context]() (xwrapcall.Observer[CtxT], ActiveFrameView) {
	o := &activeFrame[CtxT]{
		execs: []execUnit{{}},
	}
	o.execs[0].callStackInd.Store(xwrapcall.InvalidStackInd)
	return o, o
}

type ActiveFrameView interface {
	CallStackInd(execInd uint64) (xwrapcall.StackInd, bool)
}

type activeFrame[CtxT xwrapcall.Context] struct {
	execs []execUnit
}

type execUnit struct {
	callStackInd atomic.Int64
}

var (
	_ xwrapcall.Observer[xwrapcall.Context] = (*activeFrame[xwrapcall.Context])(nil)
	_ ActiveFrameView                       = (*activeFrame[xwrapcall.Context])(nil)
)

func (o *activeFrame[CtxT]) OnFrameEnter(info xwrapcall.FrameEnterEvent[CtxT]) {
	o.execs[info.ExecInd()].callStackInd.Store(info.StackIndex())
}

func (o *activeFrame[CtxT]) OnFrameLeave(info xwrapcall.FrameLeaveEvent[CtxT]) {
	o.execs[info.ExecInd()].callStackInd.Store(info.StackIndex() - 1)
}

func (o *activeFrame[CtxT]) CallStackInd(execInd uint64) (xwrapcall.StackInd, bool) {
	if execInd >= uint64(len(o.execs)) {
		return xwrapcall.InvalidStackInd, false
	}
	ind := o.execs[execInd].callStackInd.Load()
	return ind, ind != xwrapcall.InvalidStackInd
}
