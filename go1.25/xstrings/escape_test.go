package xstrings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEscape(t *testing.T) {
	testCases := []struct {
		s    string
		escC byte
		exp  string
	}{
		{
			"hello",
			'/',
			"hello",
		},
		{
			"hello / world",
			'/',
			"hello // world",
		},
		{
			"hello // world",
			'/',
			"hello //// world",
		},
		{
			"hello % world",
			'%',
			"hello %% world",
		},
		{
			"/ hello world",
			'/',
			"// hello world",
		},
		{
			"hello world /",
			'/',
			"hello world //",
		},
		{
			"hello / world",
			'o',
			"helloo / woorld",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			act := Escape(tc.s, tc.escC)
			require.Equal(t, tc.exp, act)
		})
	}
}

func TestUnescape(t *testing.T) {
	ok := false
	willPanic := true
	testCases := []struct {
		s        string
		escC     byte
		exp      string
		expPanic bool
	}{
		{
			"hello",
			'/',
			"hello",
			ok,
		},
		{
			"hello // world",
			'/',
			"hello / world",
			ok,
		},
		{
			"hello //// world",
			'/',
			"hello // world",
			ok,
		},
		{
			"hello / world",
			'/',
			"",
			willPanic,
		},
		{
			"hello %% world",
			'%',
			"hello % world",
			ok,
		},
		{
			"// hello world",
			'/',
			"/ hello world",
			ok,
		},
		{
			"hello world //",
			'/',
			"hello world /",
			ok,
		},
		{
			"helloo / woorld",
			'o',
			"hello / world",
			ok,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			if tc.expPanic {
				require.Panics(t, func() {
					Unescape(tc.s, tc.escC)
				})
			} else {
				act := Unescape(tc.s, tc.escC)
				require.Equal(t, tc.exp, act)
			}
		})
	}
}

func TestUniqueJoin(t *testing.T) {
	testCases := []struct {
		ss  []string
		sep byte
		exp string
	}{
		{
			[]string{"hello", "world"},
			'/',
			"hello/world",
		},
		{
			[]string{"hello / my ", " world"},
			'/',
			"hello // my / world",
		},
		{
			[]string{"hello ", " my / world"},
			'/',
			"hello / my // world",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.exp, func(t *testing.T) {
			act := UniqueJoin(tc.sep, tc.ss...)
			require.Equal(t, tc.exp, act)
		})
	}
}

func TestUniqueSplit(t *testing.T) {
	testCases := []struct {
		s   string
		sep byte
		exp []string
	}{
		{
			"hello/world",
			'/',
			[]string{"hello", "world"},
		},
		{
			"hello // my / world",
			'/',
			[]string{"hello / my ", " world"},
		},
		{
			"hello / my // world",
			'/',
			[]string{"hello ", " my / world"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			act := UniqueSplit(tc.s, tc.sep)
			require.Equal(t, tc.exp, act)
		})
	}
}
