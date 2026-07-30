package xsync

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/Deimvis/go-ext/go1.25/ext"
)

func WaitAll(fns ...func()) {
	wg := &sync.WaitGroup{}
	wg.Add(len(fns))
	for _, fn := range fns {
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
}

func WaitUntilFalse(
	ctx context.Context,
	fns ...func(context.Context) bool,
) {
	sfns := ext.Map(
		fns,
		func(fn func(context.Context) bool) func(context.Context, *failCountState) {
			return func(ctx context.Context, s *failCountState) {
				ok := fn(ctx)
				if !ok {
					s.fails.Add(1)
				}
			}
		},
	)
	WaitUntil(
		ctx,
		func(s *failCountState) bool {
			return s.fails.Load() > 0
		},
		sfns...,
	)
}

func WaitUntil[State any](
	ctx context.Context,
	stopFn func(*State) bool,
	fns ...func(context.Context, *State),
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var state State

	wg := &sync.WaitGroup{}
	wg.Add(len(fns))
	for _, fn := range fns {
		go func() {
			defer wg.Done()
			fn(ctx, &state)
			if stopFn(&state) {
				cancel()
			}
		}()
	}
	wg.Wait()
}

type NoState any

type failCountState struct {
	fails atomic.Uint32
}
