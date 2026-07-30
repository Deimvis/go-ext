package xwrapcall

type Observer[CtxT Context] interface {
	OnFrameEnter(FrameEnterEvent[CtxT])
	OnFrameLeave(FrameLeaveEvent[CtxT])
}

type FrameEnterEvent[CtxT Context] interface {
	ExecInd() uint64
	StackIndex() StackInd
}

type FrameLeaveEvent[CtxT Context] interface {
	FrameEnterEvent[CtxT]
	Error() error
	CalledNext() bool
	Panicked() bool
}

type frameEvent[CtxT Context] struct {
	execInd  uint64
	stackInd StackInd

	err        error
	calledNext bool
	panicked   bool
}

var (
	_ FrameEnterEvent[Context] = (*frameEvent[Context])(nil)
	_ FrameLeaveEvent[Context] = (*frameEvent[Context])(nil)
)

func (fi *frameEvent[CtxT]) ExecInd() uint64      { return fi.execInd }
func (fi *frameEvent[CtxT]) StackIndex() StackInd { return fi.stackInd }
func (fi *frameEvent[CtxT]) Error() error         { return fi.err }
func (fi *frameEvent[CtxT]) CalledNext() bool     { return fi.calledNext }
func (fi *frameEvent[CtxT]) Panicked() bool       { return fi.panicked }
