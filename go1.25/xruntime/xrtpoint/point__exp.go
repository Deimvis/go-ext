package xrtpoint

import (
	"sync"
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xruntime"
)

type Tag = string

// Point represents injected runtime point.
//
// Point is fully thread-safe.
// No change to Point is allowed after first Call() invoked.
//
// Point may be called multiple times in different code locations.
// Although, it's recommended to group points behaviour by using tags,
// not by reusing the same points.
//
// Tags help to identify points, group them and
// pass some information to callFn, because tags
// are available from Action specified in InjectionRule.
type Point interface {
	PointConst
	// WithTags clones original Point and adds new tags.
	WithTags(...Tag) Point
	Call()

	Reflect() PointReflect
}

type PointConst interface {
	Tags() []Tag
	// CallPackagePath returns package path where point's Call() was invoked.
	// It returns "" if no package information or Call() was not invoked yet.
	CallPackagePath() string
	// TODO: obfuscated package + file + line number in order to determine how many code points
	//       are affected by matching rule
	// CallLocationId() string

	ConstReflect() PointConstReflect
}

type PointReflect interface {
	PointConstReflect
}

type PointConstReflect interface {
	Injector() Injector
}

func newPoint(i Injector, callFn func(PointConst)) Point {
	return &point{
		injector: i,
		callFn:   callFn,

		tags: nil,

		cp: nil,

		freezed: &atomic.Int64{},
	}
}

type point struct {
	// immutable state
	injector Injector
	callFn   func(PointConst)

	// mutable state
	tags []Tag

	// immutable call-time state
	cp    xruntime.CallPosition
	cpSet sync.Once

	// special
	freezed *atomic.Int64
}

func (p *point) Tags() []Tag {
	return p.tags
}

func (p *point) CallPackagePath() string {
	if p.cp == nil {
		return ""
	}
	return p.cp.PackagePath()
}

func (p *point) WithTags(tags ...Tag) Point {
	p = p.clone()
	p.addTags(tags...)
	return p
}

func (p *point) Reflect() PointReflect {
	return p
}

func (p *point) ConstReflect() PointConstReflect {
	return p
}

func (p *point) Injector() Injector {
	return p.injector
}

func (p *point) addTags(tags ...Tag) {
	if p.freezed.Load() == 1 {
		panic("mutable change on freezed runtime point")
	}
	tagsRes := append(p.tags, tags...)
	ext.DeduplicateIn(&tagsRes)
	p.tags = tagsRes
}

func (p *point) Call() {
	cp := xmust.Do(xruntime.XCaller(0))
	p = p.clone()
	p.call(cp)
}

func (p *point) call(cp xruntime.CallPosition) {
	p.freezed.Store(1)
	p.cpSet.Do(func() {
		p.cp = cp
	})
	if p.cp.File() != cp.File() || p.cp.Line() != cp.Line() {
		panic("point was called from different location")
	}
	p.callFn(p)
}

func (p *point) clone() *point {
	pCopy := &point{
		injector: p.injector,
		callFn:   p.callFn,

		tags: ext.CopyElements(p.tags),

		cp: nil,

		freezed: &atomic.Int64{},
	}
	return pCopy
}
