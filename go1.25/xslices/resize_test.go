package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResize(t *testing.T) {
	tcs := []struct {
		title  string
		orig   []int
		length int
		exp    []int
	}{
		{
			"123",
			[]int{1, 2, 3},
			5,
			[]int{1, 2, 3, 0, 0},
		},
		{
			"from_empty",
			make([]int, 0, 0),
			3,
			make([]int, 3, 3),
		},
		{
			"to_empty",
			make([]int, 3, 3),
			0,
			make([]int, 0, 3),
		},
		{
			"lt_length",
			make([]int, 2, 4),
			1,
			make([]int, 1, 4),
		},
		{
			"lt_cap",
			make([]int, 2, 4),
			3,
			make([]int, 3, 4),
		},
		{
			"eq_cap",
			make([]int, 2, 4),
			4,
			make([]int, 4, 4),
		},
		{
			"gt_cap",
			make([]int, 2, 4),
			5,
			make([]int, 5, 5),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := Resize(tc.orig, tc.length)
			require.Equal(t, len(tc.exp), len(act))
			require.Equal(t, cap(tc.exp), cap(act))

			for i := range min(len(tc.orig), len(act)) {
				require.Equal(t, tc.orig[i], act[i])
			}
		})
	}
}
