package xstrconv

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatter_PerType(t *testing.T) {
	f := NewFormatter(
		FormatFns{
			PerType: TypeFormatFns{
				Int: func(i int) (string, error) {
					return strconv.Itoa(i), nil
				},
				Uint: func(ui uint) (string, error) {
					return strconv.FormatUint(uint64(ui), 10), nil
				},
				Bool: func(b bool) (string, error) {
					if b {
						return "TRUE", nil
					}
					return "FALSE", nil
				},
				Uint8: func(v byte) (string, error) {
					return string(v), nil
				},
				Int32: func(r rune) (string, error) {
					return string(r), nil
				},
				String: func(s string) (string, error) {
					return "%" + s + "%", nil
				},
			},
		},
	)

	tcs := []formatterTC{
		{
			"int",
			int(42),
			true,
			`42`,
		},
		{
			"uint",
			uint(42),
			true,
			`42`,
		},
		{
			"bool",
			bool(false),
			true,
			`FALSE`,
		},
		{
			"byte",
			byte('a'),
			true,
			`a`,
		},
		{
			"rune",
			rune('⌘'),
			true,
			`⌘`,
		},
		{
			"string",
			string("42"),
			true,
			`%42%`,
		},
		{
			"custom/string",
			MyString("42"),
			false,
			``,
		},
	}
	testFormatter(t, f, tcs)
}

func TestFormatter_PerKind(t *testing.T) {

	f := NewFormatter(
		FormatFns{
			PerKind: KindFormatFns{
				Int: func(i int) (string, error) {
					return strconv.Itoa(i), nil
				},
				Uint: func(ui uint) (string, error) {
					return strconv.FormatUint(uint64(ui), 10), nil
				},
				Bool: func(b bool) (string, error) {
					if b {
						return "TRUE", nil
					}
					return "FALSE", nil
				},
				Uint8: func(v byte) (string, error) {
					return string(v), nil
				},
				Int32: func(r rune) (string, error) {
					return string(r), nil
				},
				String: func(s string) (string, error) {
					return "%" + s + "%", nil
				},
			},
		},
	)

	tcs := []formatterTC{
		{
			"int",
			int(42),
			true,
			`42`,
		},
		{
			"uint",
			uint(42),
			true,
			`42`,
		},
		{
			"bool",
			bool(false),
			true,
			`FALSE`,
		},
		{
			"byte",
			byte('a'),
			true,
			`a`,
		},
		{
			"rune",
			rune('⌘'),
			true,
			`⌘`,
		},
		{
			"string",
			string("42"),
			true,
			`%42%`,
		},
		{
			"custom/string",
			MyString("42"),
			true,
			`%42%`,
		},
	}
	testFormatter(t, f, tcs)
}

func TestFormatter_NoFormatFn(t *testing.T) {
	f := NewFormatter(
		FormatFns{},
	)

	t.Run("int", func(t *testing.T) {
		res, err := f.Format(int(42))
		require.Error(t, err)
		require.Equal(t, ``, res)
	})
}

func TestFormatter_IntFormatPropagation(t *testing.T) {
	t.Run("int64-propagation", func(t *testing.T) {
		f := NewFormatter(
			FormatFns{
				PerType: TypeFormatFns{
					Int64: func(i int64) (string, error) {
						return strconv.FormatInt(i, 10), nil
					},
				},
			},
			WithIntFormatPropagation(),
		)

		tcs := []formatterTC{
			{
				"int64",
				int64(42),
				true,
				`42`,
			},
			{
				"int32",
				int32(42),
				true,
				`42`,
			},
			{
				"int16",
				int16(42),
				true,
				`42`,
			},
			{
				"int8",
				int8(42),
				true,
				`42`,
			},
		}
		testFormatter(t, f, tcs)
	})
	t.Run("int16-propagation", func(t *testing.T) {
		f := NewFormatter(
			FormatFns{
				PerType: TypeFormatFns{
					Int16: func(i int16) (string, error) {
						return strconv.FormatInt(int64(i), 10), nil
					},
				},
			},
			WithIntFormatPropagation(),
		)

		tcs := []formatterTC{
			{
				"int64",
				int64(42),
				false,
				``,
			},
			{
				"int32",
				int32(42),
				false,
				``,
			},
			{
				"int16",
				int16(42),
				true,
				`42`,
			},
			{
				"int8",
				int8(42),
				true,
				`42`,
			},
		}
		testFormatter(t, f, tcs)
	})
	t.Run("int-has-64-bits", func(t *testing.T) {
		defer withIntBitSize(64)()
		f := NewFormatter(
			FormatFns{
				PerType: TypeFormatFns{
					Int: func(i int) (string, error) {
						return strconv.Itoa(i), nil
					},
				},
			},
			WithIntFormatPropagation(),
		)

		tcs := []formatterTC{
			{
				"int",
				int(42),
				true,
				`42`,
			},
			{
				"int64",
				int64(42),
				true,
				`42`,
			},
			{
				"int32",
				int32(42),
				true,
				`42`,
			},
			{
				"int16",
				int16(42),
				true,
				`42`,
			},
			{
				"int8",
				int8(42),
				true,
				`42`,
			},
		}
		testFormatter(t, f, tcs)
	})
	t.Run("int-has-32-bits", func(t *testing.T) {
		defer withIntBitSize(32)()
		f := NewFormatter(
			FormatFns{
				PerType: TypeFormatFns{
					Int: func(i int) (string, error) {
						return strconv.Itoa(i), nil
					},
				},
			},
			WithIntFormatPropagation(),
		)

		tcs := []formatterTC{
			{
				"int",
				int(42),
				true,
				`42`,
			},
			{
				"int64",
				int64(42),
				false,
				``,
			},
			{
				"int32",
				int32(42),
				true,
				`42`,
			},
			{
				"int16",
				int16(42),
				true,
				`42`,
			},
			{
				"int8",
				int8(42),
				true,
				`42`,
			},
		}
		testFormatter(t, f, tcs)
	})
}

func testFormatter(t *testing.T, f Formatter, testCases []formatterTC) {
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			act, err := f.Format(tc.v)
			if tc.ok {
				require.NoError(t, err)
				require.Equal(t, tc.exp, act)
			} else {
				require.Error(t, err)
			}
		})
	}
}

type formatterTC struct {
	title string
	v     any
	ok    bool
	exp   string
}

func withIntBitSize(intBitSizeNew uintptr) func() {
	intBitSizeOld := intBitSize
	intBitSize = intBitSizeNew
	return func() {
		intBitSize = intBitSizeOld
	}
}

type MyString string
