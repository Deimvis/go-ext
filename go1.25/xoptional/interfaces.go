package xoptional

// TODO: AtomicT

type I[U any] interface {
	ConstI[U]
	ValuePtr() *U
	SetValue(U)
	Reset()
}

type ConstI[U any] interface {
	HasValue() bool
	Value() U
}
