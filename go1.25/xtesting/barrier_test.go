package xtesting

import (
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
)

func trueFn([]any) bool { return true }

func TestNewBarrier(t *testing.T) {
	parker, guard := NewBarrier()
	require.NotNil(t, parker)
	require.NotNil(t, guard)
	require.Equal(t, 0, guard.Count(trueFn))
}

func TestBarrierParker_Park(t *testing.T) {
	t.Run("unparks-when-passed", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			var returned atomic.Bool
			go func() {
				parker.Park()
				returned.Store(true)
			}()

			synctest.Wait()
			require.False(t, returned.Load(), "parker escaped")

			guard.Pass(trueFn, 1)
			synctest.Wait()
			require.True(t, returned.Load(), "parker ignored pass")
		})
	})
	t.Run("records-tags", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			go parker.Park("a", 1)

			synctest.Wait()
			hasA := func(tags []any) bool { return slices.Contains(tags, any("a")) }
			require.Equal(t, 1, guard.Count(hasA))
			require.Equal(t, 0, guard.Count(func([]any) bool { return false }))

			guard.Pass(trueFn, 1)
		})
	})
	t.Run("no-tags", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			go parker.Park()

			synctest.Wait()
			require.Equal(t, 1, guard.Count(trueFn))

			guard.Pass(trueFn, 1)
		})
	})
	t.Run("panics", func(t *testing.T) {
		t.Run("on-nil-tag", func(t *testing.T) {
			parker, _ := NewBarrier()
			require.PanicsWithError(t, "park tag <nil> of type <nil> is not comparable", func() {
				parker.Park(nil)
			})
		})
		t.Run("on-non-comparable-tag", func(t *testing.T) {
			parker, _ := NewBarrier()
			require.Panics(t, func() {
				parker.Park([]int{1, 2, 3})
			})
		})
		t.Run("no-op-when-disabled", func(t *testing.T) {
			parker, guard := NewBarrier()
			require.Equal(t, 0, guard.Disable())
			require.NotPanics(t, func() {
				parker.Park()
			})
		})
	})
}

func TestBarrierGuard_Wait(t *testing.T) {
	t.Run("blocks-until-count", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			var returned atomic.Bool
			go func() {
				guard.Wait(trueFn, 3)
				returned.Store(true)
			}()

			synctest.Wait()
			require.False(t, returned.Load())

			go parker.Park()
			go parker.Park()
			synctest.Wait()
			require.False(t, returned.Load(), "still waiting for 3rd")

			go parker.Park()
			synctest.Wait()
			require.True(t, returned.Load())

			guard.Pass(trueFn, 3)
		})
	})
	t.Run("with-pred", func(t *testing.T) {
		t.Run("matches-tags", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parker, guard := NewBarrier()

				hasTag := func(want any) func([]any) bool {
					return func(tags []any) bool { return slices.Contains(tags, want) }
				}

				var returned atomic.Bool
				go func() {
					guard.Wait(hasTag("b"), 1)
					returned.Store(true)
				}()

				go parker.Park("a")
				synctest.Wait()
				require.False(t, returned.Load(), "no matching tag yet")

				go parker.Park("b")
				synctest.Wait()
				require.True(t, returned.Load())

				guard.Pass(trueFn, 2)
			})
		})
	})
	t.Run("returns-immediately-when-already-satisfied", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			go parker.Park()
			go parker.Park()
			synctest.Wait()

			guard.Wait(trueFn, 2)

			guard.Pass(trueFn, 2)
		})
	})
	t.Run("panics-on-disable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			_, guard := NewBarrier()

			panicked := make(chan any, 1)
			go func() {
				defer func() { panicked <- recover() }()
				guard.Wait(trueFn, 5)
			}()

			synctest.Wait()
			guard.Disable()

			require.Equal(t, "barrier is disabled", <-panicked)
		})
	})
}

func TestBarrierGuard_Range(t *testing.T) {
	t.Run("visits-all-parked-in-fifo-order", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			for i := range 3 {
				go parker.Park(i)
				synctest.Wait()
			}

			var visited []int
			guard.Range(func(tags []any) bool {
				visited = append(visited, tags[0].(int))
				return true
			})
			require.Equal(t, []int{0, 1, 2}, visited)

			guard.Pass(trueFn, 3)
		})
	})
	t.Run("stops-on-false", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			for i := range 3 {
				go parker.Park(i)
				synctest.Wait()
			}

			var visited []int
			guard.Range(func(tags []any) bool {
				visited = append(visited, tags[0].(int))
				return false
			})
			require.Equal(t, []int{0}, visited)

			guard.Pass(trueFn, 3)
		})
	})
	t.Run("empty-returns-immediately", func(t *testing.T) {
		_, guard := NewBarrier()
		called := false
		guard.Range(func([]any) bool {
			called = true
			return true
		})
		require.False(t, called)
	})
}

func TestBarrierGuard_Pass(t *testing.T) {
	t.Run("releases-oldest-first", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			releasedOrder := make(chan int, 3)
			for i := range 3 {
				go func() {
					parker.Park(i)
					releasedOrder <- i
				}()
				synctest.Wait()
			}

			guard.Pass(trueFn, 1)
			synctest.Wait()
			require.Equal(t, 0, <-releasedOrder)

			guard.Pass(trueFn, 1)
			synctest.Wait()
			require.Equal(t, 1, <-releasedOrder)

			guard.Pass(trueFn, 1)
			synctest.Wait()
			require.Equal(t, 2, <-releasedOrder)
		})
	})
	t.Run("with-pred", func(t *testing.T) {
		t.Run("only-matching", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parker, guard := NewBarrier()

				released := make(chan string, 3)
				park := func(tag string) {
					go func() {
						parker.Park(tag)
						released <- tag
					}()
					synctest.Wait()
				}
				park("a")
				park("b")
				park("a")

				isA := func(tags []any) bool { return tags[0] == "a" }
				guard.Pass(isA, 2)
				synctest.Wait()

				got := []string{<-released, <-released}
				require.ElementsMatch(t, []string{"a", "a"}, got)

				require.Equal(t, 1, guard.Count(trueFn))
				guard.Pass(trueFn, 1)
			})
		})
	})
	t.Run("blocks-until-count", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			var returned atomic.Bool
			go func() {
				guard.Pass(trueFn, 2)
				returned.Store(true)
			}()

			synctest.Wait()
			require.False(t, returned.Load())

			go parker.Park()
			synctest.Wait()
			require.False(t, returned.Load())

			go parker.Park()
			synctest.Wait()
			require.True(t, returned.Load())
		})
	})
	t.Run("partial-release-leaves-rest-parked", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			for range 4 {
				go parker.Park()
				synctest.Wait()
			}

			guard.Pass(trueFn, 2)
			require.Equal(t, 2, guard.Count(trueFn))

			guard.Pass(trueFn, 2)
		})
	})
	t.Run("panics-on-disable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			_, guard := NewBarrier()

			panicked := make(chan any, 1)
			go func() {
				defer func() { panicked <- recover() }()
				guard.Pass(trueFn, 5)
			}()

			synctest.Wait()
			guard.Disable()

			require.Equal(t, "barrier is disabled", <-panicked)
		})
	})
}

func TestBarrierGuard_Disable(t *testing.T) {
	t.Run("flushes-parked-and-returns-count", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()

			var releasedCnt atomic.Int32
			for range 3 {
				go func() {
					parker.Park()
					releasedCnt.Add(1)
				}()
				synctest.Wait()
			}

			n := guard.Disable()
			synctest.Wait()

			require.Equal(t, 3, n)
			require.Equal(t, int32(3), releasedCnt.Load())
		})
	})
	t.Run("when-empty-returns-zero", func(t *testing.T) {
		_, guard := NewBarrier()
		require.Equal(t, 0, guard.Disable())
	})
	t.Run("makes-future", func(t *testing.T) {
		t.Run("park-noop", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parker, guard := NewBarrier()
				guard.Disable()

				var returned atomic.Bool
				go func() {
					parker.Park()
					returned.Store(true)
				}()

				synctest.Wait()
				require.True(t, returned.Load(), "park should return immediately when barrier is disabled")
				require.Equal(t, 0, guard.CountAll())
			})
		})
		t.Run("wait-panic", func(t *testing.T) {
			_, guard := NewBarrier()
			guard.Disable()
			require.PanicsWithValue(t, "barrier is disabled", func() {
				guard.Wait(trueFn, 1)
			})
		})
		t.Run("pass-panic", func(t *testing.T) {
			_, guard := NewBarrier()
			guard.Disable()
			require.PanicsWithValue(t, "barrier is disabled", func() {
				guard.Pass(trueFn, 1)
			})
		})
	})
	t.Run("second-call-returns-zero", func(t *testing.T) {
		_, guard := NewBarrier()
		guard.Disable()
		require.Equal(t, 0, guard.Disable())
	})
}

func TestBarrierGuard_Enable(t *testing.T) {
	t.Run("restores-parking-after-disable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			guard.Disable()
			guard.Enable()

			var returned atomic.Bool
			go func() {
				parker.Park()
				returned.Store(true)
			}()

			synctest.Wait()
			require.False(t, returned.Load(), "parker must block again after re-enable")
			require.Equal(t, 1, guard.CountAll())

			guard.Pass(trueFn, 1)
			synctest.Wait()
			require.True(t, returned.Load())
		})
	})
	t.Run("restores-wait-and-pass-after-disable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			guard.Disable()
			guard.Enable()

			require.NotPanics(t, func() {
				go parker.Park()
				synctest.Wait()
				guard.Wait(trueFn, 1)
				guard.Pass(trueFn, 1)
			})
		})
	})
	t.Run("noop-when-enabled", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			guard.Enable()

			var returned atomic.Bool
			go func() {
				parker.Park()
				returned.Store(true)
			}()

			synctest.Wait()
			require.False(t, returned.Load(), "enable on already-enabled barrier must not release parkers")

			guard.Pass(trueFn, 1)
		})
	})
}

func TestBarrierGuard_Sugar(t *testing.T) {
	t.Run("WaitOne", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			var returned atomic.Bool
			go func() {
				guard.WaitOne()
				returned.Store(true)
			}()

			synctest.Wait()
			require.False(t, returned.Load())

			go parker.Park()
			synctest.Wait()
			require.True(t, returned.Load())

			guard.Pass(trueFn, 1)
		})
	})
	t.Run("WaitN", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			var returned atomic.Bool
			go func() {
				guard.WaitN(3)
				returned.Store(true)
			}()

			for range 2 {
				go parker.Park()
			}
			synctest.Wait()
			require.False(t, returned.Load())

			go parker.Park()
			synctest.Wait()
			require.True(t, returned.Load())

			guard.Pass(trueFn, 3)
		})
	})
	t.Run("Count", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			go parker.Park("mark")
			go parker.Park("mark")
			go parker.Park("other")
			synctest.Wait()

			isMark := func(tags []any) bool { return tags[0] == "mark" }
			require.Equal(t, 2, guard.Count(isMark))

			guard.Pass(trueFn, 3)
		})
	})
	t.Run("CountHavingTags", func(t *testing.T) {
		t.Run("subset-match", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parker, guard := NewBarrier()
				go parker.Park("a", "b")
				go parker.Park("a", "c")
				synctest.Wait()

				require.Equal(t, 2, guard.CountHavingTags("a"))
				require.Equal(t, 1, guard.CountHavingTags("a", "b"))

				guard.Pass(trueFn, 2)
			})
		})
		t.Run("no-match", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parker, guard := NewBarrier()
				go parker.Park("a")
				synctest.Wait()

				require.Equal(t, 0, guard.CountHavingTags("z"))

				guard.Pass(trueFn, 1)
			})
		})
		t.Run("empty-matches-all", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				parker, guard := NewBarrier()
				go parker.Park("a")
				go parker.Park()
				synctest.Wait()

				require.Equal(t, 2, guard.CountHavingTags())

				guard.Pass(trueFn, 2)
			})
		})
	})
	t.Run("CountAll", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			for range 3 {
				go parker.Park()
			}
			synctest.Wait()

			require.Equal(t, 3, guard.CountAll())

			guard.Pass(trueFn, 3)
		})
	})
	t.Run("PassOne", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			for range 2 {
				go parker.Park()
				synctest.Wait()
			}

			guard.PassOne()
			require.Equal(t, 1, guard.CountAll())

			guard.Pass(trueFn, 1)
		})
	})
	t.Run("PassN", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			parker, guard := NewBarrier()
			for range 3 {
				go parker.Park()
				synctest.Wait()
			}

			guard.PassN(2)
			require.Equal(t, 1, guard.CountAll())

			guard.Pass(trueFn, 1)
		})
	})
}

func TestBarrier_Concurrency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		parker, guard := NewBarrier()

		const (
			kParkersPerTagCnt = 40
			kTagsCnt          = 5
			kParkersCnt       = kParkersPerTagCnt * kTagsCnt
			kRangersCnt       = 4
			kWaitersCnt       = 3
		)

		var (
			parkCallCnt     atomic.Int32
			waitCallCnt     atomic.Int32
			rangeCallCnt    atomic.Int32
			passReleasedCnt atomic.Int32
			panicCnt        atomic.Int32
		)

		var wg sync.WaitGroup

		trackPanic := func() {
			if r := recover(); r != nil {
				panicCnt.Add(1)
				t.Errorf("unexpected panic: %v", r)
			}
		}

		hasTag := func(want int) func([]any) bool {
			return func(tags []any) bool { return slices.Contains(tags, any(want)) }
		}

		rng := rand.New(rand.NewPCG(1, 2))

		for i := range kParkersCnt {
			tag := i % kTagsCnt
			wg.Go(func() {
				defer trackPanic()
				parker.Park(tag)
				parkCallCnt.Add(1)
			})
		}

		for range kRangersCnt {
			wg.Go(func() {
				defer trackPanic()
				guard.Range(func([]any) bool { return true })
				rangeCallCnt.Add(1)
			})
		}

		for range kWaitersCnt {
			tag := rng.IntN(kTagsCnt)
			wg.Go(func() {
				defer trackPanic()
				guard.Wait(hasTag(tag), 1)
				waitCallCnt.Add(1)
			})
		}

		synctest.Wait()

		for tag := range kTagsCnt {
			releaseCnt := kParkersPerTagCnt / 2
			wg.Go(func() {
				defer trackPanic()
				guard.Pass(hasTag(tag), releaseCnt)
				passReleasedCnt.Add(int32(releaseCnt))
			})
		}

		synctest.Wait()

		flushed := guard.Disable()
		wg.Wait()

		require.Equal(t, int32(0), panicCnt.Load())
		require.Equal(t, int32(kParkersCnt), parkCallCnt.Load(), "every Park must return exactly once")
		require.Equal(t, int32(kWaitersCnt), waitCallCnt.Load())
		require.Equal(t, int32(kRangersCnt), rangeCallCnt.Load())
		require.Equal(t, int32(kParkersCnt), passReleasedCnt.Load()+int32(flushed), "released + flushed should equal total parked")
	})
}

func TestValidateTag(t *testing.T) {
	type comparableStruct struct {
		A int
		B string
	}

	cases := []struct {
		title string
		tag   any
		expOk bool
	}{
		{"int", 42, true},
		{"string", "hi", true},
		{"bool", true, true},
		{"comparable-struct", comparableStruct{A: 1, B: "x"}, true},
		{"nil", nil, false},
		{"slice", []int{1, 2}, false},
		{"map", map[string]int{"a": 1}, false},
		{"func", func() {}, false},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			act := validateTag(tc.tag)
			if tc.expOk {
				require.NoError(t, act)
			} else {
				require.Error(t, act)
			}
		})
	}
}
