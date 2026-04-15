package xfn

func As[T, U any](v T) U {
	return any(v).(U)
}
