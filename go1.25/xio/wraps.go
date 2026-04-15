package xio

import "io"

// A little note for ones who interested why not straightforward approach
// for implementing io.Reader wrapper: creating struct that incapsulates io.Reader
// and overrides Read([]byte) (int, error) method.
// TL;DR - it gives no option to chain multiple stateful wrappers and promotes
// readers misuse by allowing to access wrapped io.Reader.
//
// First point there is that having a stateful wrapper the simplest way
// to expose access to the state for the user is to add another method along
// with Read([]byte) (int, error) that will serve this purpose
// (e.g. BytesRead() int64).
// But having multiple wrappers, the last wrapper will override
// all underlying wrappers' methods and user will have no access to their state.
// One solution may be for user to store a separate variable to each reader's state (wrap level),
// for example:
//   var r io.Reader
//   wr1 := NewCountingWrap(r)
//   wr2 := NewChecksumWrap(r)
//   io.ReadAll(wr2)
//   res1 := wr1.BytesRead()
//   res2 := wr2.Checksum()
// This promotes storing intermediate readers, which may be actually benefitial
// in case wrapper is needed only temporarily, but allows accidently passing
// intermediate reader forward (calling io.ReadAll(wr1) in the example above).
// We try to deal with that issue by separating Reader and wrap's state interfaces:
//   var r io.Reader
//   r, rCount := WrapReader(r, CountingReaderWrap())
//   r, rChecksum := WrapReader(r, NewChecksumReaderWrap())
//   io.ReadAll(r)
//   res1 := rCount.BytesRead()
//   res2 := rChecksum.Checksum()
// It also allows to make wrap interface a function and
// avoid creating a new struct for stateless wrappers
// (returned state is nil).

type ReaderWrapFn[State any] func(ReadFn) (ReadFn, State)

func WrapReader[State any](r io.Reader, wrapFn ReaderWrapFn[State]) (io.Reader, State) {
	fn, state := wrapFn(r.Read)
	return &readerImpl{fn: fn}, state
}

type NoState interface{}
