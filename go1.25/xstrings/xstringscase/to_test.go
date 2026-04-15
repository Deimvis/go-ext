package xstringscase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTo(t *testing.T) {
	type exp = map[Case]string
	tcs := []struct {
		title string
		text  string
		exp   exp
	}{
		{
			"empty",
			"",
			exp{
				Snake: "",
			},
		},
		{
			"simple",
			"hello world",
			exp{
				Snake: "hello_world",
			},
		},
		{
			"many spaces",
			"  hello   world   ",
			exp{
				Snake: "hello_world",
			},
		},
		{
			"camel",
			"HelloWorld",
			exp{
				Snake: "hello_world",
			},
		},
		{
			"camel",
			"HelloWorld",
			exp{
				Snake: "hello_world",
			},
		},
		{
			"kebab",
			"hello-world",
			exp{
				Snake: "hello_world",
			},
		},
		{
			"dotted",
			"hello.world",
			exp{
				Snake: "hello_world",
			},
		},
		{
			"numbers",
			"0aa-b0b-bb0",
			exp{
				Snake: "0aa_b0b_bb0",
			},
		},
		{
			"special-chars",
			"$aa-b$b-bb$",
			exp{
				Snake: "$aa_b$b_bb$",
			},
		},
		{
			"statement-ends",
			"aaa bbb \n ccc",
			exp{
				Snake: "aaa_bbb",
			},
		},
		{
			"upper-sequence",
			"aaa JSONData bbb",
			exp{
				Snake: "aaa_jsondata_bbb", // TODO: add test to check that known acronyms are preserved
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			act := ToSnake(tc.text)
			require.Equal(t, tc.exp[Snake], act)
		})
	}
}
