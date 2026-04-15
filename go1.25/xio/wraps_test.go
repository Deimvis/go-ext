package xio

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderWrapFn(t *testing.T) {
	type wrapInfo struct {
		wrapFn        ReaderWrapFn[any]
		validateState func(*testing.T, any)
	}
	testCases := []struct {
		title string
		wraps []wrapInfo
	}{
		{
			"stateless",
			[]wrapInfo{
				{
					func(readFn ReadFn) (ReadFn, any) {
						wrappedFn := func(b []byte) (int, error) {
							statelessFlag = true
							return readFn(b)
						}
						return wrappedFn, nil
					},
					func(t *testing.T, state any) {
						require.Nil(t, state)
						require.True(t, statelessFlag)
					},
				},
			},
		},

		{
			"simple_stateful",
			[]wrapInfo{
				{
					func(readFn ReadFn) (ReadFn, any) {
						s := &callCountState{}
						wrappedFn := func(b []byte) (int, error) {
							s.count++
							return readFn(b)
						}
						return wrappedFn, s
					},
					func(t *testing.T, state any) {
						s := state.(*callCountState)
						require.Equal(t, int64(2), s.CallCount())
					},
				},
			},
		},
		{
			"stateful+stateful",
			[]wrapInfo{
				{
					// call count
					func(readFn ReadFn) (ReadFn, any) {
						s := &callCountState{}
						wrappedFn := func(b []byte) (int, error) {
							s.count++
							return readFn(b)
						}
						return wrappedFn, s
					},
					func(t *testing.T, state any) {
						s := state.(*callCountState)
						require.Equal(t, int64(2), s.CallCount())
					},
				},
				{
					// byte count
					func(readFn ReadFn) (ReadFn, any) {
						s := &byteCountState{}
						wrappedFn := func(b []byte) (int, error) {
							n, err := readFn(b)
							if n > 0 {
								s.bytes += int64(n)
							}
							return n, err
						}
						return wrappedFn, s
					},
					func(t *testing.T, state any) {
						s := state.(*byteCountState)
						require.Equal(t, int64(11), s.BytesCount())
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			statelessFlag = false

			var buf bytes.Buffer
			n, err := buf.Write([]byte("hello world"))
			require.Equal(t, 11, n)
			require.NoError(t, err)

			rFn := buf.Read
			var callbacks []func(t *testing.T)
			for _, w := range tc.wraps {
				var state any
				rFn, state = w.wrapFn(rFn)
				cbFn := func(state any) func(t *testing.T) {
					return func(t *testing.T) {
						w.validateState(t, state)
					}
				}(state)
				callbacks = append(callbacks, cbFn)
			}

			chunk1 := make([]byte, 6)
			n, err = rFn(chunk1)
			require.Equal(t, 6, n)
			require.NoError(t, err)
			require.Equal(t, "hello ", string(chunk1))

			chunk2 := make([]byte, 6)
			n, err = rFn(chunk2)
			require.Equal(t, 5, n)
			require.NoError(t, err)
			require.Equal(t, "world", string(chunk2[:5]))
			require.Equal(t, byte(0), chunk2[5])

			for _, cb := range callbacks {
				cb(t)
			}
		})
	}
}

type callCountState struct {
	count int64
}

func (s *callCountState) CallCount() int64 {
	return s.count
}

type byteCountState struct {
	bytes int64
}

func (s *byteCountState) BytesCount() int64 {
	return s.bytes
}

var statelessFlag = false
