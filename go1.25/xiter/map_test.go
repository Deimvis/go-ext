package xiter

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	tcs := []struct {
		title string
		seq   []int
		fn    func(int) int
		exp   []int
	}{
		{
			"id",
			[]int{1, 2, 3},
			func(v int) int { return v },
			[]int{1, 2, 3},
		},
		{
			"empty",
			[]int{},
			func(v int) int { return v },
			nil,
		},
		{
			"x2",
			[]int{1, 2, 3},
			func(v int) int { return v * 2 },
			[]int{2, 4, 6},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := slices.Collect(Map(slices.Values(tc.seq), tc.fn))
			require.Equal(t, tc.exp, act)
		})
	}
}
