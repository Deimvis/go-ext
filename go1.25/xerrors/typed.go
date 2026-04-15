package xerrors

type TypedError[TypeT any] interface {
	error
	Type() TypeT
}
