//go:build !debug

package xrtpoint

// Inject is a shortcut to GlobalInjector.NewPoint().Call()
func Inject(opts ...InjectOption) {
}

// InjectWith is a shortcut to i.NewPoint().Call()
func InjectWith(i Injector, opts ...InjectOption) {
}

type InjectOption func(Point)

func WithTags(tags ...Tag) InjectOption {
	return func(Point) {}
}
