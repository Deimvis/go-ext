package xstringscase

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xptr"
)

func TestTokenizer(t *testing.T) {
	type pair = []int
	tcs := []struct {
		title          string
		text           string
		exp            []pair
		expEarlyStopAt *int
	}{
		{
			"empty",
			"",
			[]pair{},
			nil,
		},
		{
			"simple",
			"hello world",
			[]pair{{0, 5}, {6, 11}},
			nil,
		},
		{
			"many spaces",
			"  hello   world   ",
			[]pair{{2, 7}, {10, 15}},
			nil,
		},
		{
			"camel",
			"HelloWorld",
			[]pair{{0, 5}, {5, 10}},
			nil,
		},
		{
			"snake",
			"hello_world",
			[]pair{{0, 5}, {6, 11}},
			nil,
		},
		{
			"kebab",
			"hello-world",
			[]pair{{0, 5}, {6, 11}},
			nil,
		},
		{
			"dotted",
			"hello.world",
			[]pair{{0, 5}, {6, 11}},
			nil,
		},
		{
			"numbers",
			"0aa-b0b-bb0",
			[]pair{{0, 3}, {4, 7}, {8, 11}},
			nil,
		},
		{
			"special-chars",
			"$aa-b$b-bb$",
			[]pair{{0, 3}, {4, 7}, {8, 11}},
			nil,
		},
		{
			"statement-ends",
			"aaa bbb \n ccc",
			[]pair{{0, 3}, {4, 7}},
			xptr.T(8),
		},
		{
			"upper-sequence",
			"aaa JSONData bbb",
			[]pair{{0, 3}, {4, 12}, {13, 16}},
			nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			tzer := tokenizer{d: []byte(tc.text), i: 0}
			res := []pair{}
			for i, j, ok := tzer.Next(); ok; i, j, ok = tzer.Next() {
				res = append(res, pair{i, j})
			}
			require.Equal(t, tc.exp, res)
			if tc.expEarlyStopAt == nil {
				require.Equal(t, len(tc.text), tzer.i)
			} else {
				require.Equal(t, *tc.expEarlyStopAt, tzer.i)
			}
		})
	}
}
