package xstrconv

import (
	"fmt"
	"reflect"
)

// >>> HIGHLY EXPERIMENTAL <<<

// TODO: maybe move to xencoding/ (e.g. xrecenc - recursive encoding / xgenc - generic encoding / xmarshal). xgenc.NewMarshaler(...) encoding.BinaryMarhaler.
//       type GenericlMarhsaler interface { Marshal(MarshalContext, any) error }
//       type MarshalContext interface { Writer() bytes.Writer, ThisMarshaler() xgenc.GenericMarshaler } // errs? warnings?
//       type GenericMarshalFn[T any] func(MarshalContext, T) ([]byte, error)
//       xgenc.AsBinaryMarshaler(Marshaler) encoding.BinaryMarshaler
//       Marshal(myvar, xgenc.WithForbiddenContextUse)
//
// simple:
//
//       type MarshalFn[T any] func(T, MarshalFn[T]) ([]byte, error)
//       type GenericMarshalFn func(any) ([]byte, error)
//       ^^^ issue with allocations of []byte and possibly large errors.

// TODO: store trace to element that failed to format like valid package does? store how value matched? just more debug info

// TODO: support custom types (including interfaces: e.g. ability to use String() for those who implement fmt.Stringer)
//       note: maybe custom interface impl is same as for types: reflect.TypeFor[int](): formatFn.

// FormatFn is a function for formatting arbitrary values to string.
// It may return error since formatting may include arbitrary conversions
// and operations which aren't error-proof, unlike strconv.Format*
// functions that perform operations which are always successful.
type FormatFn[T any] func(T) (string, error)

//   - TODO: maybe rename to resolve collision with fmt.Formatter
//     (which is used for implementing custom formatting aside with type definition with Format() method).
//     This Formatter differs from fmt.Formatter that it implements independent formatting for types
//     (independent from their Format() method implementation).
//     Options for renaming so far: StandaloneFormatter, UbiquitousFormatter
type Formatter interface {
	Format(any) (string, error)
	MustFormat(any) string
}

// TypeFormatFns are format functions matched by exact type.
type TypeFormatFns struct {
	Int        FormatFn[int]
	Int8       FormatFn[int8]
	Int16      FormatFn[int16]
	Int32      FormatFn[int32] // aka rune
	Int64      FormatFn[int64]
	Uint       FormatFn[uint]
	Uint8      FormatFn[uint8] // aka byte
	Uint16     FormatFn[uint16]
	Uint32     FormatFn[uint32]
	Uint64     FormatFn[uint64]
	Uintptr    FormatFn[uintptr]
	Float32    FormatFn[float32]
	Float64    FormatFn[float64]
	Complex64  FormatFn[complex64]
	Complex128 FormatFn[complex128]
	Bool       FormatFn[bool]
	String     FormatFn[string]

	// TODO: SliceInt, SliceInt8, ...

	// TODO: support custom types
}

// KindFormatFns are format functions matched by kind (reflect.Kind).
//
// Particularly useful for auto-supporting formatting custom types
// (e.g. String will be applied `type MyString string`
// if MyString does not have type matched format function)
type KindFormatFns struct {
	Int        FormatFn[int]
	Int8       FormatFn[int8]
	Int16      FormatFn[int16]
	Int32      FormatFn[int32] // aka rune
	Int64      FormatFn[int64]
	Uint       FormatFn[uint]
	Uint8      FormatFn[uint8] // aka byte
	Uint16     FormatFn[uint16]
	Uint32     FormatFn[uint32]
	Uint64     FormatFn[uint64]
	Uintptr    FormatFn[uintptr]
	Float32    FormatFn[float32]
	Float64    FormatFn[float64]
	Complex64  FormatFn[complex64]
	Complex128 FormatFn[complex128]
	Bool       FormatFn[bool]
	String     FormatFn[string]

	// TODO: change to FormatFn[reflect.Value]
	// 	     add builders to this from simple interfaces
	Slice func() SliceFormatter
}

type FormatFns struct {
	PerType TypeFormatFns
	PerKind KindFormatFns
	// TODO: PerInterface (match when value implements interface)
}

type formatConfig struct {
	fns             FormatFns
	kindPropagation bool
}

func NewFormatter(fns FormatFns, opts ...FormatterOption) Formatter {
	cfg := formatConfig{fns: fns, kindPropagation: false}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &formatter{cfg: cfg}
}

type formatter struct {
	cfg formatConfig
}

// TODO: allow passing reflect.Value here for performance sake?
func (f *formatter) Format(v any) (string, error) {
	return f.format(v)
}

func (f *formatter) MustFormat(v any) string {
	res, err := f.Format(v)
	if err != nil {
		panic(err)
	}
	return res
}

func (f *formatter) format(v any) (string, error) {
	str, err, matched := f.formatByType(v)
	if matched {
		return str, err
	}
	str, err, matched = f.formatByKind(v)
	if matched {
		return str, err
	}
	return "", fmt.Errorf("no format function for value of type %T", v)
}

func (f *formatter) formatByType(v any) (string, error, bool) {
	str, err, matched := f.formatByTypeBuiltin(v)
	if matched {
		return str, err, true
	}
	// str, err, matched = f.formatByTypeCustom(v)
	return "", nil, false
}

func (f *formatter) formatByKind(v any) (string, error, bool) {
	fns := f.cfg.fns.PerKind
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int:
		return format(int(rv.Int()), fns.Int)
	case reflect.Int8:
		return format(int8(rv.Int()), fns.Int8)
	case reflect.Int16:
		return format(int16(rv.Int()), fns.Int16)
	case reflect.Int32:
		return format(int32(rv.Int()), fns.Int32)
	case reflect.Int64:
		return format(int64(rv.Int()), fns.Int64)
	case reflect.Uint:
		return format(uint(rv.Uint()), fns.Uint)
	case reflect.Uint8:
		return format(uint8(rv.Uint()), fns.Uint8)
	case reflect.Uint16:
		return format(uint16(rv.Uint()), fns.Uint16)
	case reflect.Uint32:
		return format(uint32(rv.Uint()), fns.Uint32)
	case reflect.Uint64:
		return format(uint64(rv.Uint()), fns.Uint64)
	case reflect.Uintptr:
		return format(uintptr(rv.Uint()), fns.Uintptr)
	case reflect.Float32:
		return format(float32(rv.Float()), fns.Float32)
	case reflect.Float64:
		return format(float64(rv.Float()), fns.Float64)
	case reflect.Complex64:
		return format(complex64(rv.Complex()), fns.Complex64)
	case reflect.Complex128:
		return format(complex128(rv.Complex()), fns.Complex128)
	case reflect.Bool:
		return format(bool(rv.Bool()), fns.Bool)
	case reflect.String:
		return format(string(rv.String()), fns.String)
	case reflect.Slice:
		sf := f.cfg.fns.PerKind.Slice()
		for i := range rv.Len() {
			elemV := rv.Index(i).Interface()
			elemVStr, err := f.format(elemV)
			if err != nil {
				return "", err, true
			}
			sf.Add(elemVStr)
		}
		str, err := sf.Format()
		return str, err, true
	default:
		return "", nil, false
	}
}

func (f *formatter) formatByTypeBuiltin(v any) (string, error, bool) {
	fns := f.cfg.fns.PerType
	switch vv := v.(type) {
	case int:
		return format(vv, fns.Int)
	case int8:
		return format(vv, fns.Int8)
	case int16:
		return format(vv, fns.Int16)
	case int32:
		return format(vv, fns.Int32)
	case int64:
		return format(vv, fns.Int64)
	case uint:
		return format(vv, fns.Uint)
	case uint8:
		return format(vv, fns.Uint8)
	case uint16:
		return format(vv, fns.Uint16)
	case uint32:
		return format(vv, fns.Uint32)
	case uint64:
		return format(vv, fns.Uint64)
	case uintptr:
		return format(vv, fns.Uintptr)
	case float32:
		return format(vv, fns.Float32)
	case float64:
		return format(vv, fns.Float64)
	case complex64:
		return format(vv, fns.Complex64)
	case complex128:
		return format(vv, fns.Complex128)
	case bool:
		return format(vv, fns.Bool)
	case string:
		return format(vv, fns.String)
	default:
		return "", nil, false
	}
}

func format[T any](v T, formatFn func(T) (string, error)) (string, error, bool) {
	if formatFn == nil {
		return "", nil, false
	}
	str, err := formatFn(v)
	return str, err, true
}
