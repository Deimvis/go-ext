package xintconv

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestShrinkToInt_32(t *testing.T) {
// 	origIntBitSize := intBitSize
// 	intBitSize = 32
// 	defer func() {
// 		intBitSize = origIntBitSize
// 	}()
// 	tcs := []struct {
// 		v   int64
// 		exp int
// 		ok  bool
// 	}{
// 		{
// 			0,
// 			0,
// 			true,
// 		},
// 		{
// 			1,
// 			1,
// 			true,
// 		},
// 		{
// 			-1,
// 			-1,
// 			true,
// 		},
// 		{
// 			int64(math.MaxInt32) + 1,
// 			0,
// 			false,
// 		},
// 		{
// 			int64(math.MinInt32) - 1,
// 			0,
// 			false,
// 		},
// 	}
// 	for _, tc := range tcs {
// 		t.Run(strconv.FormatInt(tc.v, 10), func(t *testing.T) {
// 			act, err := ShrinkToInt(tc.v)
// 			if tc.ok {
// 				require.NoError(t, err)
// 				require.Equal(t, tc.exp, act)
// 			} else {
// 				require.Error(t, err)
// 			}
// 		})
// 	}
// }

func TestShrinkToSigned(t *testing.T) {
	tcs := []struct {
		v   uint64
		exp int64
		ok  bool
	}{
		{
			0,
			0,
			true,
		},
		{
			1,
			1,
			true,
		},
		{
			uint64(math.MaxInt) + 1,
			0,
			false,
		},
	}
	for _, tc := range tcs {
		t.Run(strconv.FormatUint(tc.v, 10), func(t *testing.T) {
			act, err := ShrinkToSigned(tc.v)
			if tc.ok {
				require.NoError(t, err)
				require.Equal(t, tc.exp, act)
			} else {
				require.Error(t, err)
			}
		})
	}
}
