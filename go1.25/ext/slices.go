package ext

import (
	"cmp"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
)

// NOTE: deprecated, use xslices.Has instead
func Contains[T comparable](s []T, v T) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

// Find returns first occurance of match in a form of (value, true).
// If a match not found, second return value is false.
func Find[T any](s []T, matchFn func(v *T) bool) (T, bool) {
	for _, val := range s {
		if matchFn(&val) {
			return val, true
		}
	}
	var t T
	return t, false
}

// Map slice
func Map[T, U any](s []T, mapFn func(T) U) []U {
	res := make([]U, len(s))
	for i := range s {
		res[i] = mapFn(s[i])
	}
	return res
}

// Map slice in-place
func MapIn[T any](sp *[]T, mapFn func(T) T) []T {
	for i := range *sp {
		(*sp)[i] = mapFn((*sp)[i])
	}
	return *sp
}

// Shuffle slice
func Shuffle[T any](s []T) []T {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return ShuffleIn(&scopy)
}

// Shuffle slice in-place
func ShuffleIn[T any](sp *[]T) []T {
	s := *sp
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

// Filter filters slice.
// Keeps only those elements for which given filFn returns true.
func Filter[T any](s []T, pred func(T) bool) []T {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return FilterIn(&scopy, pred)
}

// FilterIn filters slice in-place.
// Keeps only those elements for which given filFn returns true.
func FilterIn[T any](s *[]T, pred func(T) bool) []T {
	newSz := 0
	for i := range *s {
		if pred((*s)[i]) {
			(*s)[newSz] = (*s)[i]
			newSz++
		}
	}
	*s = (*s)[:newSz]
	return *s
}

func Deduplicate[T comparable](s []T) []T {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return DeduplicateIn(&scopy)
}

// DeduplicateIn removes duplicates from slice in-place.
func DeduplicateIn[T comparable](s *[]T) []T {
	seen := make(map[T]struct{})
	return FilterIn(s, func(v T) bool {
		_, alreadyHas := seen[v]
		if alreadyHas {
			return false
		}
		seen[v] = struct{}{}
		return true
	})
}

func Sort[T any, U cmp.Ordered](s []T, keyFn func(v T) U) []T {
	scopy := make([]T, len(s))
	copy(scopy, s)
	return SortIn(&scopy, keyFn)
}

func SortIn[T any, U cmp.Ordered](s *[]T, keyFn func(v T) U) []T {
	sort.SliceStable(*s, func(i int, j int) bool {
		return keyFn((*s)[i]) < keyFn((*s)[j])
	})
	return *s
}

// https://github.com/alexanderbez/godash/blob/703c92476f3a6a947b9f2792114ecf40d7ba2c6a/godash.go#L86-L127
func SlicesEqual(slice1, slice2 interface{}) bool {
	equal, err := slicesEqual(slice1, slice2)
	if err != nil {
		panic(err)
	}
	return equal
}

func slicesEqual(slice1, slice2 interface{}) (bool, error) {
	if !IsSlice(slice1) {
		return false, fmt.Errorf("argument type '%T' is not a slice", slice1)
	} else if !IsSlice(slice2) {
		return false, fmt.Errorf("argument type '%T' is not a slice", slice2)
	}

	slice1Val := reflect.ValueOf(slice1)
	slice2Val := reflect.ValueOf(slice2)

	if slice1Val.Type().Elem() != slice2Val.Type().Elem() {
		return false, fmt.Errorf("type of '%v' does not match type of '%v'", slice1Val.Type().Elem(), slice2Val.Type().Elem())
	}

	if slice1Val.Len() != slice2Val.Len() {
		return false, nil
	}

	result := true
	i, n := 0, slice1Val.Len()

	for i < n {
		j := 0
		e := false
		for j < n && !e {
			if slice1Val.Index(i).Interface() == slice2Val.Index(j).Interface() {
				e = true
			}
			j++
		}
		if !e {
			result = false
		}
		i++
	}

	return result, nil
}

func IsSlice(value interface{}) bool {
	kind := reflect.ValueOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func CopyElements[T any](s []T) []T {
	sCopy := make([]T, len(s))
	copy(sCopy, s)
	return sCopy
}

func ReferenceElements[T any](s []T) []*T {
	return Map(s, func(v T) *T { return &v })
}

func DereferenceElements[T any](s []*T) []T {
	return Map(s, func(v *T) T { return *v })
}
