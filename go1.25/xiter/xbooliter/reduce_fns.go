package xbooliter

import (
	"github.com/Deimvis/go-ext/go1.25/xiter"
)

var All xiter.ReduceFn[bool] = func(cur bool, v bool) bool {
	return cur && v
}
