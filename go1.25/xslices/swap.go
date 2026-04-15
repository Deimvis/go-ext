package xslices

type SwapFn func(i, j int)

func NewSwapFn[T any, S ~[]T](s S) SwapFn {
	return func(i, j int) {
		s[i], s[j] = s[j], s[i]
	}
}
