//go:build debug

package xrtpoint

import (
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xruntime"
)

// Inject is a shortcut to GlobalInjector.NewPoint().Call()
func Inject(opts ...InjectOption) {
	cp := xmust.Do(xruntime.XCaller(1))
	p := GlobalInjector.NewPoint()
	for _, opt := range opts {
		p = opt(p)
	}
	pp := p.(*point)
	pp.call(cp)
}

// InjectWith is a shortcut to i.NewPoint().Call()
func InjectWith(i Injector, opts ...InjectOption) {
	cp := xmust.Do(xruntime.XCaller(1))
	p := i.NewPoint()
	for _, opt := range opts {
		p = opt(p)
	}
	pp := p.(*point)
	pp.call(cp)
}

type InjectOption func(Point) Point

func WithTags(tags ...Tag) InjectOption {
	return func(p Point) Point {
		return p.WithTags(tags...)
	}
}
