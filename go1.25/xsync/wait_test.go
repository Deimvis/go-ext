package xsync

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitAll(t *testing.T) {
	t.Run("runs_all-functions", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 16
			var counter atomic.Int32
			fns := make([]func(), n)
			for i := range fns {
				fns[i] = func() { counter.Add(1) }
			}

			WaitAll(fns...)

			require.Equal(t, int32(n), counter.Load())
		})
	})

	t.Run("runs_in-parallel", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 8
			var started atomic.Int32
			release := make(chan struct{})

			fns := make([]func(), n)
			for i := range fns {
				fns[i] = func() {
					started.Add(1)
					<-release
				}
			}

			go WaitAll(fns...)

			synctest.Wait()
			require.Equal(t, int32(n), started.Load(), "all goroutines should have started before any can return")

			close(release)
		})
	})

	t.Run("zero-functions", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			WaitAll()
		})
	})

	t.Run("waits_for-slowest", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var slowDone atomic.Bool

			fast := func() {}
			slow := func() {
				time.Sleep(50 * time.Millisecond)
				slowDone.Store(true)
			}

			WaitAll(fast, slow)
			require.True(t, slowDone.Load(), "WaitAll must return only after the slowest fn completed")
		})
	})
}

func TestUntil(t *testing.T) {
	t.Run("stops_when-stopfn-true", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const blockers = 5

			var trigger, blockerExited atomic.Int32
			fns := []func(context.Context, *NoState){
				func(ctx context.Context, _ *NoState) {
					trigger.Add(1)
				},
			}
			for range blockers {
				fns = append(fns, func(ctx context.Context, _ *NoState) {
					<-ctx.Done()
					blockerExited.Add(1)
				})
			}

			Until(
				context.Background(),
				func(*NoState) bool { return trigger.Load() >= 1 },
				fns...,
			)

			require.Equal(t, int32(blockers), blockerExited.Load(), "all blockers should exit via ctx cancel triggered by stopFn")
		})
	})

	t.Run("all-run_when-stopfn-never-true", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 12

			var ranCount atomic.Int32
			fns := make([]func(context.Context, *NoState), n)
			for i := range fns {
				fns[i] = func(ctx context.Context, _ *NoState) {
					ranCount.Add(1)
				}
			}

			Until(
				context.Background(),
				func(*NoState) bool { return false },
				fns...,
			)

			require.Equal(t, int32(n), ranCount.Load())
		})
	})

	t.Run("zero-functions", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			Until(context.Background(), func(*NoState) bool { return false })
		})
	})

	t.Run("parent-ctx-cancellation_propagates", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 6

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			var ranCount atomic.Int32
			fns := make([]func(context.Context, *NoState), n)
			for i := range fns {
				fns[i] = func(ctx context.Context, _ *NoState) {
					<-ctx.Done()
					ranCount.Add(1)
				}
			}

			Until(ctx, func(*NoState) bool { return false }, fns...)

			require.Equal(t, int32(n), ranCount.Load())
		})
	})

	t.Run("stopfn-invoked_after-each-fn", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 7

			var ranCount, stopCallCount atomic.Int32
			fns := make([]func(context.Context, *NoState), n)
			for i := range fns {
				fns[i] = func(ctx context.Context, _ *NoState) {
					ranCount.Add(1)
				}
			}

			Until(
				context.Background(),
				func(*NoState) bool {
					stopCallCount.Add(1)
					return false
				},
				fns...,
			)

			require.Equal(t, int32(n), ranCount.Load())
			require.Equal(t, int32(n), stopCallCount.Load())
		})
	})
}

func TestUntilFalse(t *testing.T) {
	t.Run("stops_on-first-false", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 8

			var exited atomic.Int32
			fns := make([]func(context.Context) bool, n)
			fns[0] = func(ctx context.Context) bool {
				exited.Add(1)
				return false
			}
			for i := 1; i < n; i++ {
				fns[i] = func(ctx context.Context) bool {
					<-ctx.Done()
					exited.Add(1)
					return true
				}
			}

			UntilFalse(context.Background(), fns...)

			require.Equal(t, int32(n), exited.Load())
		})
	})

	t.Run("all-true_all-run", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 10

			var ran atomic.Int32
			fns := make([]func(context.Context) bool, n)
			for i := range fns {
				fns[i] = func(ctx context.Context) bool {
					ran.Add(1)
					return true
				}
			}

			UntilFalse(context.Background(), fns...)

			require.Equal(t, int32(n), ran.Load())
		})
	})

	t.Run("multiple-falses_are-safe", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 10
			const falses = 4

			var exited atomic.Int32
			fns := make([]func(context.Context) bool, n)
			for i := range falses {
				fns[i] = func(ctx context.Context) bool {
					exited.Add(1)
					return false
				}
			}
			for i := falses; i < n; i++ {
				fns[i] = func(ctx context.Context) bool {
					<-ctx.Done()
					exited.Add(1)
					return true
				}
			}

			UntilFalse(context.Background(), fns...)

			require.Equal(t, int32(n), exited.Load())
		})
	})

	t.Run("zero-functions", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			UntilFalse(context.Background())
		})
	})

	t.Run("respects_parent-ctx", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			const n = 6

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			var exited atomic.Int32
			fns := make([]func(context.Context) bool, n)
			for i := range fns {
				fns[i] = func(ctx context.Context) bool {
					<-ctx.Done()
					exited.Add(1)
					return true
				}
			}

			UntilFalse(ctx, fns...)

			require.Equal(t, int32(n), exited.Load())
		})
	})
}
