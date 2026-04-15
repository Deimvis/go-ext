package ext

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	testCases := []findTC[int]{
		{
			[]int{1, 2, 3},
			func(x *int) bool { return (*x)%2 != 0 },
			ptr(1),
		},
		{
			[]int{2, 3},
			func(x *int) bool { return (*x)%2 != 0 },
			ptr(3),
		},
		{
			[]int{2, 4, 6},
			func(x *int) bool { return (*x)%2 != 0 },
			nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			sCopy := append([]int(nil), tc.s...)
			actual, ok := Find(tc.s, tc.matchFn)
			if tc.expected != nil {
				require.Equal(t, true, ok)
				require.Equal(t, *tc.expected, actual)
			} else {
				require.Equal(t, false, ok)
			}
			require.Equal(t, tc.s, sCopy) // not changed
		})
	}
}

func TestFilter(t *testing.T) {
	testCases := []filterTC[int]{
		{
			[]int{1, 2, 3},
			func(x int) bool { return x%2 != 0 },
			[]int{1, 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			initialCopy := append([]int(nil), tc.initial...)
			actual := Filter(tc.initial, tc.filFn)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, initialCopy, tc.initial) // not changed
		})
	}
}

func TestFilterIn(t *testing.T) {
	testCases := []filterTC[int]{
		{
			[]int{1, 2, 3},
			func(x int) bool { return x%2 != 0 },
			[]int{1, 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := FilterIn(&tc.initial, tc.filFn)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, tc.expected, tc.initial) // changed
		})
	}
}

func TestShuffleIn(t *testing.T) {
	testCases := []struct {
		initial []int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			initialCopy := append([]int(nil), tc.initial...)
			actual := ShuffleIn(&tc.initial)
			require.Equal(t, len(initialCopy), len(actual))
			require.NotEqual(t, initialCopy, actual)
			require.Equal(t, tc.initial, actual) // changed

		})
	}
}

func TestDeduplicateIn(t *testing.T) {
	testCases := []struct {
		title    string
		initial  []int
		expected []int // sorted
	}{
		{
			"simple",
			[]int{0, 0, 1, 1, 2, 2, 3, 3, 3, 3},
			[]int{0, 1, 2, 3},
		},
		{
			"empty slice",
			[]int{},
			[]int{},
		},
		{
			"already deduplicated",
			[]int{0, 1, 2, 3},
			[]int{0, 1, 2, 3},
		},
		{
			"mixed duplicates",
			[]int{0, 1, 0, 2, 1, 3, 3, 2},
			[]int{0, 1, 2, 3},
		},
		{
			"one duplicate",
			[]int{0, 0, 0, 0, 0},
			[]int{0},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.title), func(t *testing.T) {
			actual := DeduplicateIn(&tc.initial)
			slices.Sort(actual)
			require.Equal(t, tc.expected, actual)

		})
	}
}

func TestSortIn(t *testing.T) {
	keyFn := func(x int) int { return x }
	testCases := []struct {
		title    string
		initial  []int
		expected []int
	}{
		{
			"simple",
			[]int{0, 2, 1, 3},
			[]int{0, 1, 2, 3},
		},
		{
			"empty slice",
			[]int{},
			[]int{},
		},
		{
			"already sorted",
			[]int{0, 1, 2, 3},
			[]int{0, 1, 2, 3},
		},
		{
			"mixed",
			[]int{3, 0, 1, 10, 9},
			[]int{0, 1, 3, 9, 10},
		},
		{
			"with duplicates",
			[]int{1, 5, 0, 0, 4, 5},
			[]int{0, 0, 1, 4, 5, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := SortIn(&tc.initial, keyFn)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, tc.initial, actual) // in-place
		})
	}
}

func TestSortIn_KeyFn(t *testing.T) {
	type st struct {
		i int
		s string
	}

	t.Run("int", func(t *testing.T) {
		keyFnInt := func(v st) int { return v.i }
		testCasesInt := []struct {
			title    string
			initial  []st
			expected []st
			keyFn    func(v st) int
		}{
			{
				"simple",
				[]st{{i: 1}, {i: 2}, {i: 0}},
				[]st{{i: 0}, {i: 1}, {i: 2}},
				keyFnInt,
			},
			{
				"empty slice",
				[]st{},
				[]st{},
				keyFnInt,
			},
			{
				"with duplicates",
				[]st{{i: 3}, {i: 0}, {i: 0}, {i: 5}, {i: 4}, {i: 1}, {i: 2}, {i: 4}},
				[]st{{i: 0}, {i: 0}, {i: 1}, {i: 2}, {i: 3}, {i: 4}, {i: 4}, {i: 5}},
				keyFnInt,
			},
		}
		for _, tc := range testCasesInt {
			t.Run(tc.title, func(t *testing.T) {
				actual := SortIn(&tc.initial, tc.keyFn)
				require.Equal(t, tc.expected, actual)
				require.Equal(t, tc.initial, actual) // in-place
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		keyFnString := func(v st) string { return v.s }
		testCasesString := []struct {
			title    string
			initial  []st
			expected []st
			keyFn    func(v st) string
		}{
			{
				"simple",
				[]st{{s: "c"}, {s: "a"}, {s: "b"}},
				[]st{{s: "a"}, {s: "b"}, {s: "c"}},
				keyFnString,
			},
			{
				"empty slice",
				[]st{},
				[]st{},
				keyFnString,
			},
			{
				"with duplicates",
				[]st{{s: "c"}, {s: "a"}, {s: "a"}, {s: "f"}, {s: "c"}, {s: "d"}, {s: "b"}},
				[]st{{s: "a"}, {s: "a"}, {s: "b"}, {s: "c"}, {s: "c"}, {s: "d"}, {s: "f"}},
				keyFnString,
			},
		}
		for _, tc := range testCasesString {
			t.Run(tc.title, func(t *testing.T) {
				actual := SortIn(&tc.initial, tc.keyFn)
				require.Equal(t, tc.expected, actual)
				require.Equal(t, tc.initial, actual) // in-place
			})
		}
	})
}

type findTC[T any] struct {
	s        []T
	matchFn  func(*T) bool
	expected *T
}

type filterTC[T any] struct {
	initial  []T
	filFn    func(T) bool
	expected []T
}
