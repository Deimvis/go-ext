package xptr

func ShallowClone[T any](p *T) *T {
	v := *p
	return &v
}
