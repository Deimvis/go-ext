package xwrapcallctx

import "context"

type copying[Self any] interface {
	StdContext() context.Context
	CopyOnto(Context) Self
}

type WithCopying[Self any] interface {
	Context
	copying[Self]
}
