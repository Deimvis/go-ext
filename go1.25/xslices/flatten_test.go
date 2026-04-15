package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlatten(t *testing.T) {
	tcs := []struct {
		title string
		inp   [][]int
		exp   []int
	}{
		{
			"empty",
			[][]int{},
			[]int{},
		},
		{
			"123",
			[][]int{{1, 2, 3}},
			[]int{1, 2, 3},
		},
		{
			"123,",
			[][]int{{1, 2, 3}, {}},
			[]int{1, 2, 3},
		},
		{
			",123",
			[][]int{{}, {1, 2, 3}},
			[]int{1, 2, 3},
		},
		{
			",123,",
			[][]int{{}, {1, 2, 3}, {}},
			[]int{1, 2, 3},
		},
		{
			"1,2,3",
			[][]int{{1}, {2}, {3}},
			[]int{1, 2, 3},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := Flatten(tc.inp)
			require.Equal(t, tc.exp, act)
			require.Equal(t, len(act), cap(act))
		})
	}
}
