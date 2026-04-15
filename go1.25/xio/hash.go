package xio

// TODO: NewHashReaderWrap(hashType ...)

// something like this:

// type md55ReaderWrapState struct {
// 	hash hash.Hash
// }

// func (s md55ReaderWrapState) Sum() []byte {
// 	return s.hash.Sum(nil)
// }

// func md5ReaderWrap(fn xio.ReadFn) (xio.ReadFn, md55ReaderWrapState) {
// 	var state md55ReaderWrapState
// 	state.hash = md5.New()
// 	wrappedFn := func(p []byte) (int, error) {
// 		n, err := fn(p)
// 		if n > 0 {
// 			nn, suberr := state.hash.Write(p[:n])
// 			xmust.NoErr(suberr)
// 			xmust.Eq(nn, n)
// 		}
// 		return n, err
// 	}
// 	return wrappedFn, state
// }
