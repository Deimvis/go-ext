package xwrapcallobs

import (
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
)

func NewCallStackInd[CtxT xwrapcall.Context]() (xwrapcall.Observer[CtxT], CallStackIndView) {
	o := &callStackInd[CtxT]{
		execs: []execUnit{{}},
	}
	o.execs[0].activeStackInd.Store(xwrapcall.InvalidStackInd)
	return o, o
}

type CallStackIndView interface {
	ActiveStackIndex(execInd uint64) (xwrapcall.StackInd, bool)
}

type callStackInd[CtxT xwrapcall.Context] struct {
	execs []execUnit
}

type execUnit struct {
	activeStackInd atomic.Int64
}

var (
	_ xwrapcall.Observer[xwrapcall.Context] = (*callStackInd[xwrapcall.Context])(nil)
	_ CallStackIndView                      = (*callStackInd[xwrapcall.Context])(nil)
)

func (o *callStackInd[CtxT]) OnFrameEnter(info xwrapcall.FrameEnterEvent[CtxT]) {
	o.execs[info.ExecInd()].activeStackInd.Store(info.StackIndex())
}

func (o *callStackInd[CtxT]) OnFrameLeave(info xwrapcall.FrameLeaveEvent[CtxT]) {
	o.execs[info.ExecInd()].activeStackInd.Store(info.StackIndex() - 1)
}

func (o *callStackInd[CtxT]) ActiveStackIndex(execInd uint64) (xwrapcall.StackInd, bool) {
	if execInd >= uint64(len(o.execs)) {
		return xwrapcall.InvalidStackInd, false
	}
	ind := o.execs[execInd].activeStackInd.Load()
	return ind, ind != xwrapcall.InvalidStackInd
}
