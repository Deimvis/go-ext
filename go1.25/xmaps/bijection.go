package xmaps

import "fmt"

func ValidateBijection[K comparable, V comparable, M ~map[K]V](m M) error {
	seen := make(map[V]struct{}, len(m)) // optimistic allocation
	for _, v := range m {
		if _, ok := seen[v]; ok {
			return fmt.Errorf("duplicate value: %v", v)
		}
		seen[v] = struct{}{}
	}
	return nil
}
