package ext

func IfElse[T any](cond bool, v1 T, v2 T) T {
	if cond {
		return v1
	}
	return v2
}
