package xmaps

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateBijection(t *testing.T) {
	type m = map[string]int
	testCases := []struct {
		title string
		orig  m
		exp   error
	}{
		{
			"simple",
			m{"a": 1, "b": 2},
			nil,
		},
		{
			"empty",
			m{},
			nil,
		},
		{
			"many_keys",
			m{"a": 1, "c": 2, "y": 3, "x": 4},
			nil,
		},
		{
			"not_bijection",
			m{"a": 1, "c": 2, "y": 3, "x": 4, "z": 1},
			errors.New("duplicate value: 1"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			origCopy := copyMap(tc.orig)
			act := ValidateBijection(tc.orig)
			require.Equal(t, origCopy, tc.orig)
			require.Equal(t, tc.exp, act)
		})
	}
}
