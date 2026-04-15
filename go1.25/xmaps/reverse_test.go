package xmaps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverse(t *testing.T) {
	type m = map[string]int
	type revm = map[int]string
	testCases := []struct {
		title string
		orig  m
		exp   revm
	}{
		{
			"simple",
			m{"a": 1, "b": 2},
			revm{1: "a", 2: "b"},
		},
		{
			"empty",
			m{},
			revm{},
		},
		{
			"many_keys",
			m{"a": 1, "c": 2, "y": 3, "x": 4},
			revm{1: "a", 2: "c", 3: "y", 4: "x"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			origCopy := copyMap(tc.orig)
			act := Reverse(tc.orig)
			require.Equal(t, origCopy, tc.orig)
			require.Equal(t, tc.exp, act)
		})
	}
}

func copyMap[K comparable, V any, M ~map[K]V](m M) M {
	res := make(M, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}
