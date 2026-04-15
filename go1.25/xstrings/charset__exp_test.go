package xstrings

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/runes"
)

func TestBelongsTo(t *testing.T) {
	tcs := []struct {
		title string
		s     string
		rs    runes.Set
		exp   bool
	}{
		{
			"abc",
			"abc",
			runes.Predicate(func(r rune) bool {
				return slices.Contains([]rune{'a', 'b', 'c'}, r)
			}),
			true,
		},
		{
			"always false",
			"arbitrary string абв 世界",
			runes.Predicate(func(r rune) bool {
				return false
			}),
			false,
		},
		{
			"always true",
			"arbitrary string абв 世界",
			runes.Predicate(func(r rune) bool {
				return true
			}),
			true,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := BelongsTo(tc.s, tc.rs)
			require.Equal(t, tc.exp, act)
		})
	}
}
