package xtesting

type _GateParker interface {
	Park(tags ...any)
}

type _GateGuard interface {
	Open()
	Close()

	WaitPassed(pred func(tags []any) bool, count int)
	WaitPassedN(count int)
}
