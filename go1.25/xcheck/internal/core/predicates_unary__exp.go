package core

// TODO: replace with interface everywhere predicate is used
type UnaryPredicate[T any] func(T) bool
