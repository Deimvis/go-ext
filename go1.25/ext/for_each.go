package ext

import "github.com/Deimvis/go-ext/go1.25/xcheck/xmust"

type runForEachSettings struct {
	parallel bool
	untilFn  func(any, error) bool // finished when untilFn returns true
}

type RunForEachOption func(*runForEachSettings)

func WithSequential() RunForEachOption {
	return func(opts *runForEachSettings) {
		opts.parallel = false
	}
}

func WithUntilFn(fn func(any, error) bool) RunForEachOption {
	return func(opts *runForEachSettings) {
		opts.untilFn = fn
	}
}

// TODO: move to xfn
func RunForEach[T any, U any](fn func(T) (U, error), args []T, opts ...RunForEachOption) {
	settings := runForEachSettings{
		parallel: false,
		untilFn:  nil,
	}
	for _, opt := range opts {
		opt(&settings)
	}

	xmust.True(settings.parallel == false, "parallel is not supported yet")
	var v U
	var err error
	for _, arg := range args {
		v, err = fn(arg)
		if settings.untilFn != nil && settings.untilFn(v, err) {
			break
		}
	}
}

// TODO: just apply function to each element of slice should be xslices.ForEach + xslices.ForEach (parallel version)
// this is needed to run function over parameters with until function
// the issue is whether it should return value on last function call or return nothing
// or it should provide return callback with options, like WithLastReturnValue(&value)
