package xwrapcallobs

import (
	"slices"
	"sync"
	"time"

	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xtime"
)

// TODO: add option to pass storage to timings
// which will be responsible for size boundaries,
// eviction, downsampling, etc.

func NewTimings[CtxT xwrapcall.Context](clock xtime.Clock) (xwrapcall.Observer[CtxT], TimingsView) {
	t := &timings[CtxT]{clock: clock}
	return t, t
}

type TimingsView interface {
	Events() []TimingsEvent
}

type TimingsEvent struct {
	Ts         time.Time
	Kind       TimingsEventKind
	ExecInd    uint64
	StackIndex xwrapcall.StackInd
}

type TimingsEventKind int

const (
	FrameEnter TimingsEventKind = iota
	FrameLeave
)

type timings[CtxT xwrapcall.Context] struct {
	clock xtime.Clock

	mu     sync.Mutex
	events []TimingsEvent
}

var (
	_ xwrapcall.Observer[xwrapcall.Context] = (*timings[xwrapcall.Context])(nil)
	_ TimingsView                           = (*timings[xwrapcall.Context])(nil)
)

func (t *timings[CtxT]) OnFrameEnter(info xwrapcall.FrameEnterEvent[CtxT]) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, TimingsEvent{
		Ts:         t.clock.Now(),
		Kind:       FrameEnter,
		ExecInd:    info.ExecInd(),
		StackIndex: info.StackIndex(),
	})
}

func (t *timings[CtxT]) OnFrameLeave(info xwrapcall.FrameLeaveEvent[CtxT]) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, TimingsEvent{
		Ts:         t.clock.Now(),
		Kind:       FrameLeave,
		ExecInd:    info.ExecInd(),
		StackIndex: info.StackIndex(),
	})
}

func (t *timings[CtxT]) Events() []TimingsEvent {
	t.mu.Lock()
	defer t.mu.Unlock()
	return slices.Clone(t.events)
}
