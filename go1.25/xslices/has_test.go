package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHas(t *testing.T) {
	tcs := []struct {
		title string
		s     []int
		v     int
		exp   bool
	}{
		{
			"123/true",
			[]int{1, 2, 3},
			2,
			true,
		},
		{
			"123/false",
			[]int{1, 2, 3},
			4,
			false,
		},
		{
			"empty",
			[]int{},
			2,
			false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := Has(tc.s, tc.v)
			require.Equal(t, tc.exp, act)
		})
	}
}
