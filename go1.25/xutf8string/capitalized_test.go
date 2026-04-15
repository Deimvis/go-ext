package xutf8string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCapitalized(t *testing.T) {
	testCases := []isCapitalizedTestCase{
		{
			"hello world - true",
			"Hello, World!",
			true,
		},
		{
			"hello world - false",
			"hello, World!",
			false,
		},
		{
			"one letter - true",
			"A",
			true,
		},
		{
			"one letter - false",
			"a",
			false,
		},
		{
			"empty string - false",
			"",
			false,
		},
		{
			"not unicode - true",
			"Б",
			true,
		},
		{
			"not unicode - false",
			"б",
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%d", tc.title, i), func(t *testing.T) {
			actual := IsCapitalized(tc.s)
			require.Equal(t, tc.expected, actual)
		})
	}
}

type isCapitalizedTestCase struct {
	title    string
	s        string
	expected bool
}
