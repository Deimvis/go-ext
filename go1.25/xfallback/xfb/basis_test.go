package xfb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOn(t *testing.T) {
	tcs := []onTC[int]{
		{
			"trivial_never",
			func(int) bool { return false },
			0,
			666,
			0,
		},
		{
			"trivial_always",
			func(int) bool { return true },
			0,
			666,
			666,
		},
		{
			"eq0",
			func(v int) bool { return v == 0 },
			0,
			666,
			666,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := On(tc.pred, tc.v, tc.fb)
			require.Equal(t, tc.exp, act)
		})
	}
}

type onTC[T any] struct {
	title string
	pred  func(T) bool
	v     T
	fb    T
	exp   T
}
