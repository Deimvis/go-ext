package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xmaps"
)

func TestToInclusionMap(t *testing.T) {
	tcs := []struct {
		title string
		s     []int
		exp   map[int]struct{}
	}{
		{
			"123",
			[]int{1, 2, 3},
			map[int]struct{}{
				1: {},
				2: {},
				3: {},
			},
		},
		{
			"empty",
			[]int{},
			map[int]struct{}{},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := ToInclusionMap(tc.s)
			require.Equal(t, xmaps.InclusionMap[int](tc.exp), act)
		})
	}
}
