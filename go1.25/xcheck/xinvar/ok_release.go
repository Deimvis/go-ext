//go:build !debug

package xinvar

func Ok[T any](v T, ok bool, msgAndArgs ...interface{}) T {
	return v
}
