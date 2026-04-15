package xwrapcall

type AbortOption func(AbortInfoMutable)

func WithReason(r string) AbortOption {
	return func(i AbortInfoMutable) {
		i.SetReason(r)
	}
}

func WithFields(fields ...Field) AbortOption {
	return func(i AbortInfoMutable) {
		i.SetFields(fields...)
	}
}

func WithAutoFields(keysAndValues ...any) AbortOption {
	panic("not implemented")
}

func WithCallStack() AbortOption {
	panic("not implemented")
}
