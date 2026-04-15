package balg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCNF(t *testing.T) {
	testCases := []struct {
		title string
		e     Expression
		exp   Expression
	}{
		{
			"literal",
			L_TRUE,
			L_TRUE,
		},
		{
			"variable",
			Var("a"),
			Var("a"),
		},
		{
			"and-literals",
			And(L_TRUE, L_FALSE),
			And(L_TRUE, L_FALSE),
		},
		{
			"and-variables",
			And(Var("a"), Var("b")),
			And(Var("a"), Var("b")),
		},
		{
			"and-literals_n_variables",
			And(L_TRUE, Var("a")),
			And(L_TRUE, Var("a")),
		},
		{
			"or-literals",
			Or(L_TRUE, L_FALSE),
			Or(L_TRUE, L_FALSE),
		},
		{
			"or-variables",
			Or(Var("a"), Var("b")),
			Or(Var("a"), Var("b")),
		},
		{
			"or-literals_n_variables",
			Or(L_TRUE, Var("a")),
			Or(L_TRUE, Var("a")),
		},
		{
			"and-or-literals",
			And(Or(L_TRUE, L_FALSE), Or(L_FALSE, L_FALSE)),
			And(Or(L_TRUE, L_FALSE), Or(L_FALSE, L_FALSE)),
		},
		{
			"and-or-and-variables",
			And(
				Or(
					Var("a"),
					Var("b"),
					And(Var("c"), Var("d")),
				),
			),
			And(
				Or(Var("a"), Var("b"), Var("c")),
				Or(Var("a"), Var("b"), Var("d")),
			),
		},
		{
			"or-and-variables",
			Or(
				Var("a"),
				And(
					Var("b"),
					Var("c"),
				),
			),
			And(
				Or(
					Var("a"),
					Var("b"),
				),
				Or(
					Var("a"),
					Var("c"),
				),
			),
		},
		{
			"or-and-or-variables",
			// https://www.wolframalpha.com/input?i=%28A+or+%28B+and+C%29+or+%28D+and+%28E+or+F%29%29%29
			Or(
				Var("a"),
				And(
					Var("b"),
					Var("c"),
				),
				And(
					Var("d"),
					Or(
						Var("e"),
						Var("f"),
					),
				),
			),
			And(
				Or(
					Var("a"),
					Var("b"),
					Var("d"),
				),
				Or(
					Var("a"),
					Var("b"),
					Var("e"),
					Var("f"),
				),
				Or(
					Var("a"),
					Var("c"),
					Var("d"),
				),
				Or(
					Var("a"),
					Var("c"),
					Var("e"),
					Var("f"),
				),
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			act := ToCNF(tc.e)
			require.Equal(t, ExpressionToString(tc.exp), ExpressionToString(act))
		})
	}
}
