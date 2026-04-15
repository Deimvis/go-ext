package xslices

func Has[T comparable, S ~[]T](s S, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}
