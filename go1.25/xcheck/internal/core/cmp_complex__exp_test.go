package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnyOf(t *testing.T) {
	tcs := []struct {
		title      string
		pred       BinaryPredicate_[int]
		needle     int
		haystack   []int
		msgAndArgs []any
		exp        bool
		expMsg     string
	}{
		{
			"123",
			EqPred[int]{},
			2,
			[]int{1, 2, 3},
			[]any{},
			true,
			"",
		},
		{
			"empty",
			EqPred[int]{},
			1,
			[]int{},
			[]any{PrintValues()},
			false,
			"any_of not met: 1 != <x>, where <x> belongs to []",
		},
		{
			"no value",
			EqPred[int]{},
			2,
			[]int{1, 3},
			[]any{PrintValues()},
			false,
			"any_of not met: 2 != <x>, where <x> belongs to [1 3]",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act, actMsg := AnyOf(tc.pred, tc.needle, tc.haystack, tc.msgAndArgs...)
			require.Equal(t, tc.exp, act)
			require.Equal(t, tc.expMsg, actMsg)
		})
	}
}
