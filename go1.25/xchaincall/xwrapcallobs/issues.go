package xwrapcallobs

import (
	"slices"
	"sync"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
)

func NewIssues[CtxT xwrapcall.Context]() (xwrapcall.Observer[CtxT], IssuesView) {
	i := &issues[CtxT]{}
	return i, i
}

type IssuesView interface {
	Events() []IssueEvent
}

type IssueEvent struct {
	ExecInd    uint64
	StackIndex xwrapcall.StackInd
	Err        error
	Panicked   bool
}

type issues[CtxT xwrapcall.Context] struct {
	mu     sync.Mutex
	events []IssueEvent
}

var (
	_ xwrapcall.Observer[xwrapcall.Context] = (*issues[xwrapcall.Context])(nil)
	_ IssuesView                            = (*issues[xwrapcall.Context])(nil)
)

func (i *issues[CtxT]) OnFrameEnter(xwrapcall.FrameEnterEvent[CtxT]) {}

func (i *issues[CtxT]) OnFrameLeave(info xwrapcall.FrameLeaveEvent[CtxT]) {
	if !info.Panicked() && info.Error() == nil {
		return
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.events = append(i.events, IssueEvent{
		ExecInd:    info.ExecInd(),
		StackIndex: info.StackIndex(),
		Err:        info.Error(),
		Panicked:   info.Panicked(),
	})
}

func (i *issues[CtxT]) Events() []IssueEvent {
	i.mu.Lock()
	defer i.mu.Unlock()
	return slices.Clone(i.events)
}
