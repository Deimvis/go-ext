package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFill(t *testing.T) {
	tcs := []struct {
		title string
		s     []int
		newFn func(int) int
		exp   []int
	}{
		{
			"empty",
			nil,
			nil,
			nil,
		},
		{
			"000->111",
			[]int{0, 0, 0},
			func(int) int { return 1 },
			[]int{1, 1, 1},
		},
		{
			"000->123",
			[]int{0, 0, 0},
			func(i int) int { return i + 1 },
			[]int{1, 2, 3},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := tc.s
			Fill(&act, tc.newFn)
			require.Equal(t, tc.exp, act)
		})
	}
}

func TestFillConst(t *testing.T) {
	tcs := []struct {
		title string
		s     []int
		v     int
		exp   []int
	}{
		{
			"empty",
			nil,
			0,
			nil,
		},
		{
			"000->111",
			[]int{0, 0, 0},
			1,
			[]int{1, 1, 1},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := tc.s
			FillConst(&act, tc.v)
			require.Equal(t, tc.exp, act)
		})
	}
}
