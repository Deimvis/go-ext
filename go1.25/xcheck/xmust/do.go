package xmust

// Do panics if err != nil and does nothing otherwise.
// Format string from msgAndArgs must contain last %s for error.
func Do[T any](v T, err error, msgAndArgs ...interface{}) T {
	NoErr(err, msgAndArgs...)
	return v
}

// TODO: decide whether to keep Do0

// Do0 panics if err != nil and does nothing otherwise.
// Format string from msgAndArgs must contain last %s for error.
// func Do0(err error, msgAndArgs ...interface{}) {
// 	NoErr(err, msgAndArgs...)
// }

// Do1 panics if err != nil and does nothing otherwise.
// Format string from msgAndArgs must contain last %s for error.
func Do1[T any](v T, err error, msgAndArgs ...interface{}) T {
	NoErr(err, msgAndArgs...)
	return v
}

// Do2 panics if err != nil and does nothing otherwise.
// Format string from msgAndArgs must contain last %s for error.
func Do2[T, U any](v1 T, v2 U, err error, msgAndArgs ...interface{}) (T, U) {
	NoErr(err, msgAndArgs...)
	return v1, v2
}

// CheckDo panics if any err != nil and does nothing otherwise.
func CheckDo(ret ...interface{}) {
	for _, v := range ret {
		err, ok := v.(error)
		if ok && err != nil {
			panic(err)
		}
	}
}
