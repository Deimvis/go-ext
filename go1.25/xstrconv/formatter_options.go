package xstrconv

import (
	"unsafe"
)

type FormatterOption func(*formatConfig)

// TODO: WithFloatFormatPropagation
// TODO: WithFloatOptimalFormatFn(significantDigits int)  // (finds best precision for each value separately)

func WithIntFormatPropagation() FormatterOption {
	propagate := func(
		intFormat *FormatFn[int],
		int64Format *FormatFn[int64],
		int32Format *FormatFn[int32],
		int16Format *FormatFn[int16],
		int8Format *FormatFn[int8],
	) {
		if (*intFormat) == nil {
			if intBitSize <= 64 && (*int64Format) != nil {
				*intFormat = func(v int) (string, error) {
					return (*int64Format)(int64(v))
				}
			}
		}

		if (*int64Format) == nil {
			if intBitSize >= 64 && (*intFormat) != nil {
				*int64Format = func(v int64) (string, error) {
					return (*intFormat)(int(v))
				}
			}
		}

		if (*int32Format) == nil {
			if intBitSize >= 32 && (*intFormat) != nil {
				*int32Format = func(v int32) (string, error) {
					return (*intFormat)(int(v))
				}
			} else if (*int64Format) != nil {
				*int32Format = func(v int32) (string, error) {
					return (*int64Format)(int64(v))
				}
			}
		}

		if (*int16Format) == nil {
			if intBitSize >= 16 && (*intFormat) != nil {
				*int16Format = func(v int16) (string, error) {
					return (*intFormat)(int(v))
				}
			} else if (*int64Format) != nil {
				*int16Format = func(v int16) (string, error) {
					return (*intFormat)(int(v))
				}
			} else if (*int32Format) != nil {
				*int16Format = func(v int16) (string, error) {
					return (*int32Format)(int32(v))
				}
			}
		}

		if (*int8Format) == nil {
			if intBitSize >= 8 && (*intFormat) != nil {
				*int8Format = func(v int8) (string, error) {
					return (*intFormat)(int(v))
				}
			} else if (*int64Format) != nil {
				*int8Format = func(v int8) (string, error) {
					return (*int64Format)(int64(v))
				}
			} else if (*int32Format) != nil {
				*int8Format = func(v int8) (string, error) {
					return (*int32Format)(int32(v))
				}
			} else if (*int16Format) != nil {
				*int8Format = func(v int8) (string, error) {
					return (*int16Format)(int16(v))
				}
			}
		}
	}
	return func(cfg *formatConfig) {
		propagate(
			&cfg.fns.PerType.Int,
			&cfg.fns.PerType.Int64,
			&cfg.fns.PerType.Int32,
			&cfg.fns.PerType.Int16,
			&cfg.fns.PerType.Int8,
		)

		propagate(
			&cfg.fns.PerKind.Int,
			&cfg.fns.PerKind.Int64,
			&cfg.fns.PerKind.Int32,
			&cfg.fns.PerKind.Int16,
			&cfg.fns.PerKind.Int8,
		)
	}
}

func WithUintFormatPropagation() FormatterOption {
	propagate := func(
		uintFormat *FormatFn[uint],
		uintptrFormat *FormatFn[uintptr],
		uint64Format *FormatFn[uint64],
		uint32Format *FormatFn[uint32],
		uint16Format *FormatFn[uint16],
		uint8Format *FormatFn[uint8],
	) {
		if (*uintFormat) == nil {
			if uintBitSize <= 64 && (*uint64Format) != nil {
				*uintFormat = func(v uint) (string, error) {
					return (*uint64Format)(uint64(v))
				}
			}
		}

		if (*uintptrFormat) == nil {
			if uintptrBitSize <= 64 && (*uint64Format) != nil {
				*uintptrFormat = func(v uintptr) (string, error) {
					return (*uint64Format)(uint64(v))
				}
			}
		}

		if (*uint64Format) == nil {
			if uintBitSize >= 64 && (*uintFormat) != nil {
				*uint64Format = func(v uint64) (string, error) {
					return (*uintFormat)(uint(v))
				}
			} else if uintptrBitSize >= 64 && (*uintptrFormat) != nil {
				*uint64Format = func(v uint64) (string, error) {
					return (*uintptrFormat)(uintptr(v))
				}
			}
		}

		if (*uint32Format) == nil {
			if uintBitSize >= 32 && (*uintFormat) != nil {
				*uint32Format = func(v uint32) (string, error) {
					return (*uintFormat)(uint(v))
				}
			} else if uintptrBitSize >= 32 && (*uintptrFormat) != nil {
				*uint32Format = func(v uint32) (string, error) {
					return (*uintptrFormat)(uintptr(v))
				}
			} else if (*uint64Format) != nil {
				*uint32Format = func(v uint32) (string, error) {
					return (*uint64Format)(uint64(v))
				}
			}
		}

		if (*uint16Format) == nil {
			if uintBitSize >= 16 && (*uintFormat) != nil {
				*uint16Format = func(v uint16) (string, error) {
					return (*uintFormat)(uint(v))
				}
			} else if uintptrBitSize >= 16 && (*uintptrFormat) != nil {
				*uint16Format = func(v uint16) (string, error) {
					return (*uintptrFormat)(uintptr(v))
				}
			} else if (*uint64Format) != nil {
				*uint16Format = func(v uint16) (string, error) {
					return (*uint64Format)(uint64(v))
				}
			} else if (*uint32Format) != nil {
				*uint16Format = func(v uint16) (string, error) {
					return (*uint32Format)(uint32(v))
				}
			}
		}

		if (*uint8Format) == nil {
			if uintBitSize >= 8 && (*uintFormat) != nil {
				*uint8Format = func(v uint8) (string, error) {
					return (*uintFormat)(uint(v))
				}
			} else if uintptrBitSize >= 8 && (*uintptrFormat) != nil {
				*uint8Format = func(v uint8) (string, error) {
					return (*uintptrFormat)(uintptr(v))
				}
			} else if (*uint64Format) != nil {
				*uint8Format = func(v uint8) (string, error) {
					return (*uint64Format)(uint64(v))
				}
			} else if (*uint32Format) != nil {
				*uint8Format = func(v uint8) (string, error) {
					return (*uint32Format)(uint32(v))
				}
			} else if (*uint16Format) != nil {
				*uint8Format = func(v uint8) (string, error) {
					return (*uint16Format)(uint16(v))
				}
			}
		}
	}
	return func(cfg *formatConfig) {
		propagate(
			&cfg.fns.PerType.Uint,
			&cfg.fns.PerType.Uintptr,
			&cfg.fns.PerType.Uint64,
			&cfg.fns.PerType.Uint32,
			&cfg.fns.PerType.Uint16,
			&cfg.fns.PerType.Uint8,
		)

		propagate(
			&cfg.fns.PerKind.Uint,
			&cfg.fns.PerKind.Uintptr,
			&cfg.fns.PerKind.Uint64,
			&cfg.fns.PerKind.Uint32,
			&cfg.fns.PerKind.Uint16,
			&cfg.fns.PerKind.Uint8,
		)
	}
}

func WithByteFormat(fn FormatFn[byte]) FormatterOption {
	return func(cfg *formatConfig) {
		cfg.fns.PerType.Uint8 = fn
	}
}

func WithRuneFormat(fn FormatFn[rune]) FormatterOption {
	return func(cfg *formatConfig) {
		cfg.fns.PerType.Int32 = fn
	}
}

var (
	intBitSize     = unsafe.Sizeof(int(0)) * 8
	uintBitSize    = unsafe.Sizeof(uint(0)) * 8
	uintptrBitSize = unsafe.Sizeof(uintptr(0)) * 8
)
