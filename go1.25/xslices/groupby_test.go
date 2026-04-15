package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupBy(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		type is = []int
		keyFn := func(x int) int { return x }
		testCases := []groupByTestCase[int, int]{
			{
				"unique",
				[]int{1, 2, 3},
				keyFn,
				map[int][]int{1: is{1}, 2: is{2}, 3: is{3}},
			},
			{
				"same",
				[]int{1, 1, 1},
				keyFn,
				map[int][]int{1: is{1, 1, 1}},
			},
			{
				"random",
				[]int{1, 3, 2, 2, 4, 5, 1, 6},
				keyFn,
				map[int][]int{1: is{1, 1}, 2: is{2, 2}, 3: is{3}, 4: is{4}, 5: is{5}, 6: is{6}},
			},
		}
		runGroupByTests(t, testCases)
	})
	t.Run("struct", func(t *testing.T) {
		type As = []A
		keyFn := func(a A) string { return a.Value }
		testCases := []groupByTestCase[A, string]{
			{
				"unique",
				[]A{{"a"}, {"b"}, {"c"}},
				keyFn,
				map[string][]A{"a": As{{"a"}}, "b": As{{"b"}}, "c": As{{"c"}}},
			},
			{
				"same",
				[]A{{"a"}, {"a"}, {"a"}},
				keyFn,
				map[string][]A{"a": As{{"a"}, {"a"}, {"a"}}},
			},
			{
				"random",
				[]A{{"a"}, {"c"}, {"b"}, {"b"}, {"d"}, {"e"}, {"a"}, {"f"}},
				keyFn,
				map[string][]A{"a": As{{"a"}, {"a"}}, "b": As{{"b"}, {"b"}}, "c": As{{"c"}}, "d": As{{"d"}}, "e": As{{"e"}}, "f": As{{"f"}}},
			},
		}
		runGroupByTests(t, testCases)
	})
}

func TestGroupIndBy(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		type inds = []int
		keyFn := func(x int) int { return x }
		testCases := []groupIndByTestCase[int, int]{
			{
				"unique",
				[]int{1, 2, 3},
				keyFn,
				map[int][]int{1: inds{0}, 2: inds{1}, 3: inds{2}},
			},
			{
				"same",
				[]int{1, 1, 1},
				keyFn,
				map[int][]int{1: inds{0, 1, 2}},
			},
			{
				"random",
				[]int{1, 3, 2, 2, 4, 5, 1, 6},
				keyFn,
				map[int][]int{1: inds{0, 6}, 2: inds{2, 3}, 3: inds{1}, 4: inds{4}, 5: inds{5}, 6: inds{7}},
			},
		}
		runGroupIndByTests(t, testCases)
	})
	t.Run("struct", func(t *testing.T) {
		type inds = []int
		keyFn := func(a A) string { return a.Value }
		testCases := []groupIndByTestCase[A, string]{
			{
				"unique",
				[]A{{"a"}, {"b"}, {"c"}},
				keyFn,
				map[string][]int{"a": inds{0}, "b": inds{1}, "c": inds{2}},
			},
			{
				"same",
				[]A{{"a"}, {"a"}, {"a"}},
				keyFn,
				map[string][]int{"a": inds{0, 1, 2}},
			},
			{
				"random",
				[]A{{"a"}, {"c"}, {"b"}, {"b"}, {"d"}, {"e"}, {"a"}, {"f"}},
				keyFn,
				map[string][]int{"a": inds{0, 6}, "b": inds{2, 3}, "c": inds{1}, "d": inds{4}, "e": inds{5}, "f": inds{7}},
			},
		}
		runGroupIndByTests(t, testCases)
	})
}

func runGroupByTests[T any, K comparable](t *testing.T, testCases []groupByTestCase[T, K]) {
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			act := GroupBy(tc.inp, tc.keyFn)
			require.Equal(t, tc.exp, act)
		})
	}
}

func runGroupIndByTests[T any, K comparable](t *testing.T, testCases []groupIndByTestCase[T, K]) {
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			act := GroupIndBy(tc.inp, tc.keyFn)
			require.Equal(t, tc.exp, act)
		})
	}
}

type groupByTestCase[T any, K comparable] struct {
	title string
	inp   []T
	keyFn func(T) K
	exp   map[K][]T
}

type groupIndByTestCase[T any, K comparable] struct {
	title string
	inp   []T
	keyFn func(T) K
	exp   map[K][]int
}

type A struct {
	Value string
}
