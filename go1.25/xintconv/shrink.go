package xintconv

import (
	"fmt"
	"math"
)

func ShrinkToInt(v int64) (int, error) {
	if intBitSize < 64 {
		if v < math.MinInt {
			return 0, fmt.Errorf("min int overflow (%d)", v)
		}
		if v > math.MaxInt {
			return 0, fmt.Errorf("max int overflow (%d)", v)
		}
	}
	return int(v), nil
}

func ShrinkToSigned(v uint64) (int64, error) {
	if v > uint64(math.MaxInt64) {
		return 0, fmt.Errorf("max int overflow (%d)", v)
	}
	return int64(v), nil
}

const (
	intBitSize  = 32 << (^uint(0) >> 63) // from "math" package
	uintBitSize = intBitSize             // according to "math" package
)
