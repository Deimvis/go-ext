package xcontext

import "context"

// Map is a proxy to context
// with map semantics.
// Map is NOT thread-safe.
type Map interface {
	Get(key any) (any, bool)
	Set(key any, value any)
	Context() context.Context
}

func NewMap(c context.Context) Map {
	return &ctxMap{c: c}
}

type ctxMap struct {
	c context.Context
}

var _ Map = &ctxMap{}

func (cm ctxMap) Get(key any) (any, bool) {
	v := cm.c.Value(key)
	return v, v != nil
}

func (cm *ctxMap) Set(key any, value any) {
	cm.c = context.WithValue(cm.c, key, value)
}

func (cm ctxMap) Context() context.Context {
	return cm.c
}
