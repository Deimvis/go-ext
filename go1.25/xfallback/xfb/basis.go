package xfb

// GOEXPERIMENT=aliastypeparams
// type OnFn[T any] = func(T) bool
// type FallbackFn[T any] = func(v T, fallback T) T

func On[T any](onFn func(T) bool, v T, fallback T) T {
	if !onFn(v) {
		return v
	}
	return fallback
}

func New[T any](onFn func(T) bool) func(T, T) T {
	return func(v T, fallback T) T {
		return On(onFn, v, fallback)
	}
}
