package xchans

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: go1.24, consider using synctest: synctest.Wait after launching run goroutine in concurrent tests

func TestMapToSlice(t *testing.T) {
	type dst = []int
	type exp = []int
	tcs := []struct {
		title       string
		runWithChan func(run func(chan int))
		fn          func(int) int
		dst         []int
		exp         []int
	}{
		{
			"123",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
				run(ch)
			},
			id,
			dst{},
			exp{1, 2, 3},
		},
		{
			"123 override",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
				run(ch)
			},
			id,
			dst{777, 888, 999}, // some trash here
			exp{1, 2, 3},
		},
		{
			"no elements",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				close(ch)
				run(ch)
			},
			id,
			dst{},
			exp{},
		},
		{
			"fn x2",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
				run(ch)
			},
			func(v int) int { return v * 2 },
			dst{},
			exp{2, 4, 6},
		},
		{
			"concurrent",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				done := make(chan struct{})

				ch <- 1
				go func() {
					run(ch)
					done <- struct{}{}
				}()
				// synctest.Wait()
				ch <- 2
				ch <- 3

				close(ch)
				<-done
			},
			id,
			dst{},
			exp{1, 2, 3},
		},
		{
			"liveness",
			func(run func(chan int)) {
				ch := make(chan int, 1)
				done := make(chan struct{})

				go func() {
					ch <- 1
					ch <- 2
					ch <- 3
					close(ch)
					done <- struct{}{}
				}()
				run(ch)

				<-done
			},
			id,
			dst{},
			exp{1, 2, 3},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			var act []int
			var ch_ chan int
			run := func(ch chan int) {
				ch_ = ch
				act = MapToSlice(ch, tc.fn, tc.dst)
			}
			tc.runWithChan(run)
			require.Equal(t, tc.exp, act)
			_, ok := <-ch_
			require.False(t, ok)
		})
	}
}

func TestMapNToSlice(t *testing.T) {
	type dst = []int
	type exp = []int
	type surplus = []int
	tcs := []struct {
		title          string
		runWithChan    func(run func(chan int))
		fn             func(int) int
		dst            []int
		n              int
		exp            []int
		expChanSurplus []int
	}{
		{
			"123",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
				run(ch)
			},
			id,
			dst{},
			3,
			exp{1, 2, 3},
			surplus{},
		},
		{
			"123 override",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
				run(ch)
			},
			id,
			dst{777, 888, 999}, // some trash here
			3,
			exp{1, 2, 3},
			surplus{},
		},
		{
			"no elements",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				close(ch)
				run(ch)
			},
			id,
			dst{},
			3,
			exp{},
			surplus{},
		},
		{
			"fn x2",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				close(ch)
				run(ch)
			},
			func(v int) int { return v * 2 },
			dst{},
			3,
			exp{2, 4, 6},
			surplus{},
		},
		{
			"concurrent",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				done := make(chan struct{})

				ch <- 1
				go func() {
					run(ch)
					done <- struct{}{}
				}()
				ch <- 2
				ch <- 3

				close(ch)
				<-done
			},
			id,
			dst{},
			3,
			exp{1, 2, 3},
			surplus{},
		},
		{
			"liveness",
			func(run func(chan int)) {
				ch := make(chan int, 1)
				done := make(chan struct{})

				go func() {
					ch <- 1
					ch <- 2
					ch <- 3
					close(ch)
					done <- struct{}{}
				}()
				run(ch)

				<-done
			},
			id,
			dst{},
			3,
			exp{1, 2, 3},
			surplus{},
		},
		{
			"with surplus",
			func(run func(chan int)) {
				ch := make(chan int, 5)
				ch <- 1
				ch <- 2
				ch <- 3
				ch <- 4
				ch <- 5
				close(ch)
				run(ch)
			},
			id,
			dst{},
			3,
			exp{1, 2, 3},
			surplus{4, 5},
		},
		{
			"concurrent with surplus",
			func(run func(chan int)) {
				ch := make(chan int, 2)
				done := make(chan struct{})

				go func() {
					ch <- 1
					ch <- 2
					ch <- 3
					ch <- 4
					ch <- 5
					close(ch)
					done <- struct{}{}
				}()
				run(ch)

				<-done
			},
			id,
			dst{},
			3,
			exp{1, 2, 3},
			surplus{4, 5},
		},
		{
			"closed before N",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				ch <- 1
				close(ch)
				run(ch)
			},
			id,
			dst{},
			3,
			exp{1},
			surplus{},
		},
		{
			"concurrent closed before N",
			func(run func(chan int)) {
				ch := make(chan int, 3)
				done := make(chan struct{})

				go func() {
					ch <- 1
					close(ch)
					done <- struct{}{}
				}()
				run(ch)
			},
			id,
			dst{},
			3,
			exp{1},
			surplus{},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			var act []int
			var ch_ chan int
			run := func(ch chan int) {
				ch_ = ch
				act = MapNToSlice(ch, tc.fn, tc.dst, tc.n)
			}
			tc.runWithChan(run)
			require.Equal(t, tc.exp, act)

			actSurplus := []int{}
			for {
				v, ok := <-ch_
				if !ok {
					break
				}
				actSurplus = append(actSurplus, v)
			}
			require.Equal(t, tc.expChanSurplus, actSurplus)
		})
	}
}

func id(v int) int {
	return v
}
