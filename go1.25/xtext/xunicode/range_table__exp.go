package xunicode

import "unicode"

func NewRangeTable(opts ...rangeTableOpt) *unicode.RangeTable {
	panic("not implemented")
	// TODO:
	// - apply options
	// - sort ranges and calc MaxLatin for ranges16
}

type rangeTableOpt = func(rt *unicode.RangeTable)

func WithSequentialRange(first rune, last rune) rangeTableOpt {
	panic("not implemented")
	// if last <= 0xFFFF {
	//  // add range16
	// 	lo := uint16(frst)
	//  hi := uint16(last)
	// } else {
	//  // add range32
	// }
}
