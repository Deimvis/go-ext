package xwrapcallctx

type AbortOption func(AbortInfoMutable)

func AbortWithReason(r string) AbortOption {
	return func(i AbortInfoMutable) {
		i.SetReason(r)
	}
}

func AbortWithFields(fields ...Field) AbortOption {
	return func(i AbortInfoMutable) {
		i.SetFields(fields...)
	}
}

func AbortWithAutoFields(keysAndValues ...any) AbortOption {
	panic("not implemented")
}

func AbortWithCallStack() AbortOption {
	panic("not implemented")
}
