package xhash

import (
	"fmt"
	"hash"
)

func Calc(h hash.Hash, p []byte) []byte {
	n, err := h.Write(p)
	if err != nil {
		panic(fmt.Errorf("unexpectedly got error on hash Write() call: %w", err))
	}
	if n != len(p) {
		panic(fmt.Errorf("unexpectedly hash Write() call did not consume full data: %d != %d", n, len(p)))
	}
	return h.Sum(nil)
}
