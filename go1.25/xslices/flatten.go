package xslices

func Flatten[T any, S ~[]T](s []S) S {
	sz := 0
	for i := range s {
		sz += len(s[i])
	}
	res := make(S, sz)
	offset := 0
	for i := range s {
		for j := range s[i] {
			res[offset] = s[i][j]
			offset++
		}
	}
	return res
}
