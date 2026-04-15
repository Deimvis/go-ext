//go:build !debug

package xrtpoint

var (
	noopInector Injector = &noopInjectorImpl{}
)

type noopInjectorImpl struct{}

func (i noopInjectorImpl) NewPoint() Point {
	return newPoint(i, func(PointConst) {})
}

func (i noopInjectorImpl) Rules() OrderedInjectionRules {
	panic(errUseInRelease)
}
