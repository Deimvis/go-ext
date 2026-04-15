package xiter

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	tcs := []struct {
		title string
		seq   []int
		fn    func(int, int) int
		init  int
		exp   int
	}{
		{
			"sum",
			[]int{1, 2, 3},
			func(cur int, v int) int {
				return cur + v
			},
			0,
			6,
		},
		{
			"empty",
			[]int{},
			func(cur int, v int) int {
				return cur + v
			},
			123,
			123,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := Reduce(slices.Values(tc.seq), tc.fn, tc.init)
			require.Equal(t, tc.exp, act)
		})
	}
}
