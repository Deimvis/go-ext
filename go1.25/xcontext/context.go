package xcontext

import (
	"context"
	"time"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

func WithTimeoutIn(ctx *context.Context, timeout time.Duration) context.CancelFunc {
	xmust.NotNilPtr(ctx)
	newCtx, cancel := context.WithTimeout(*ctx, timeout)
	*ctx = newCtx
	return cancel
}
