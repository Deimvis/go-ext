package xmaps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyMap(t *testing.T) {
	testCases := []struct {
		origMap map[string]int
		check   func(t *testing.T, origMap map[string]int, copyMap map[string]int)
	}{
		{
			make(map[string]int),
			func(t *testing.T, origMap map[string]int, copyMap map[string]int) {
				require.Equal(t, 0, len(origMap))
				require.Equal(t, 0, len(copyMap))
				copyMap["key"] = 42
				require.Equal(t, 0, len(origMap))
				require.Equal(t, 1, len(copyMap))
			},
		},
		{
			map[string]int{
				"key":   42,
				"other": 0,
			},
			func(t *testing.T, origMap map[string]int, copyMap map[string]int) {
				require.Equal(t, 2, len(origMap))
				require.Equal(t, 2, len(copyMap))
				require.Equal(t, 42, origMap["key"])
				require.Equal(t, 42, copyMap["key"])
				copyMap["key"] = 0
				require.Equal(t, 42, origMap["key"])
				require.Equal(t, 0, copyMap["key"])
				origMap["key"] = 40
				require.Equal(t, 40, origMap["key"])
				require.Equal(t, 0, copyMap["key"])

				delete(copyMap, "key")
				require.Equal(t, 2, len(origMap))
				require.Equal(t, 1, len(copyMap))
				require.Equal(t, 40, origMap["key"])
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			origMapCopy := Clone(tc.origMap)
			tc.check(t, tc.origMap, origMapCopy)
		})
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		origMap  map[string]int
		filFn    func(k string, v int) bool
		expected map[string]int
		inplace  bool
	}{
		{
			make(map[string]int),
			func(k string, v int) bool { return true },
			make(map[string]int),
			true,
		},
		{
			map[string]int{"key0": 0, "key1": 1, "key2": 2},
			func(k string, v int) bool { return v%2 == 0 },
			map[string]int{"key0": 0, "key2": 2},
			true,
		},
		{
			map[string]int{"key0": 0, "key1": 1, "key2": 2},
			func(k string, v int) bool { return v%2 == 0 },
			map[string]int{"key0": 0, "key2": 2},
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var res map[string]int
			if tc.inplace {
				FilterIn(&tc.origMap, tc.filFn)
				res = tc.origMap
			} else {
				origMapCopy := Clone(tc.origMap)
				res = Filter(tc.origMap, tc.filFn)
				require.Equal(t, origMapCopy, tc.origMap)
			}
			require.Equal(t, tc.expected, res)
		})
	}
}

func TestMapKeys(t *testing.T) {
	testCases := []struct {
		m        map[string]int
		expected []string
	}{
		{
			make(map[string]int),
			nil,
		},
		{
			map[string]int{},
			[]string{},
		},
		{
			map[string]int{"a": 1},
			[]string{"a"},
		},
		{
			map[string]int{"a": 1, "b": 1},
			[]string{"a", "b"},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.ElementsMatch(t, tc.expected, Keys(tc.m))
		})
	}
}

func TestHasKey(t *testing.T) {
	testCases := []struct {
		m        map[string]int
		key      string
		expected bool
	}{
		{
			make(map[string]int),
			"something",
			false,
		},
		{
			map[string]int{"a": 1, "b": 2},
			"a",
			true,
		},
		{
			map[string]int{"a": 1, "b": 2},
			"b",
			true,
		},
		{
			map[string]int{"a": 1, "b": 2},
			"c",
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.Equal(t, tc.expected, HasKey(tc.m, tc.key))
		})
	}
}

func TestMerge(t *testing.T) {
	testCases := []struct {
		left     map[string]int
		right    map[string]int
		expected map[string]int
	}{
		{
			make(map[string]int),
			make(map[string]int),
			make(map[string]int),
		},
		{
			map[string]int{"a": 1},
			map[string]int{"b": 2},
			map[string]int{"a": 1, "b": 2},
		},
		{
			map[string]int{"a": 1},
			map[string]int{"a": 2},
			map[string]int{"a": 2},
		},
		{
			map[string]int{"a": 1, "b": 2},
			map[string]int{"a": 2, "c": 3},
			map[string]int{"a": 2, "b": 2, "c": 3},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			leftcopy := Clone(tc.left)
			rightcopy := Clone(tc.right)
			res := Merge(tc.left, tc.right)
			require.Equal(t, leftcopy, tc.left)
			require.Equal(t, rightcopy, tc.right)
			require.Equal(t, tc.expected, res)
		})
	}
}
