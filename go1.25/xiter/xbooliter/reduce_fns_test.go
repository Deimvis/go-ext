package xbooliter

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xiter"
)

func TestAll(t *testing.T) {
	tcs := []struct {
		title string
		seq   []bool
		init  bool
		exp   bool
	}{
		{
			"all-true",
			[]bool{true, true, true},
			true,
			true,
		},
		{
			"empty",
			[]bool{},
			true,
			true,
		},
		{
			"one-false",
			[]bool{true, false, true},
			true,
			false,
		},
		{
			"all-false",
			[]bool{false, false, false},
			true,
			false,
		},
		{
			"all-true/init-false",
			[]bool{true, true, true},
			false,
			false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := xiter.Reduce(slices.Values(tc.seq), All, tc.init)
			require.Equal(t, tc.exp, act)
		})
	}
}
