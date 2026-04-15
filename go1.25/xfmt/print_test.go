package xfmt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSprintfg(t *testing.T) {
	type args = []any
	testCases := []struct {
		msgAndArgs args
		expected   string
	}{
		{
			args{},
			"",
		},
		{
			args{""},
			"",
		},
		{
			args{"kek"},
			"kek",
		},
		{
			args{"%s", "kek"},
			"kek",
		},
		{
			args{"%d + %s = ERROR", 2, "2"},
			"2 + 2 = ERROR",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := Sprintfg(tc.msgAndArgs...)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintkv(t *testing.T) {
	type args = []any
	testCases := []struct {
		title         string
		kvFormat      string
		sep           string
		keysAndValues []any
		expected      string
		ok            bool
	}{
		{
			"empty",
			"",
			"",
			args{},
			"",
			true,
		},
		{
			"a=b",
			"%v=%v",
			",",
			args{"a", "b", 1, 2},
			`a=b,1=2`,
			true,
		},
		{
			"plain types",
			"%v=%v",
			" ",
			args{int(1), int64(2), uint(3), uint64(4), float32(4.5), float64(5.5), "a", "b"},
			`1=2 3=4 4.5=5.5 a=b`,
			true,
		},
		{
			"dangling key",
			"",
			"",
			args{1, 2, 3},
			``,
			false,
		},
		{
			"key-value formatting string",
			"%s=%d",
			"",
			args{"key", 42},
			`key=42`,
			true,
		},
		{
			"bad formatting string",
			"%d=%s",
			"",
			args{"key", 42},
			`%!d(string=key)=%!s(int=42)`,
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			if tc.ok {
				actual := Sprintfkv(tc.kvFormat, tc.sep, tc.keysAndValues...)
				require.Equal(t, tc.expected, actual)
			} else {
				require.Panics(t, func() {
					Sprintfkv(tc.kvFormat, tc.sep, tc.keysAndValues...)
				})
			}
		})
	}
}
