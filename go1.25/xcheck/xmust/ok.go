package xmust

func Ok[T any](v T, ok bool, msgAndArgs ...interface{}) T {
	True(ok, msgAndArgs...)
	return v
}

func Ok0(ok bool, msgAndArgs ...interface{}) {
	True(ok, msgAndArgs...)
}

func Ok2[T1, T2 any](v1 T1, v2 T2, ok bool, msgAndArgs ...interface{}) (T1, T2) {
	True(ok, msgAndArgs...)
	return v1, v2
}
