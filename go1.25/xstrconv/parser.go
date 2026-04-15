package xstrconv

// type ParseFn[T any] func(string, *T) error

// type Parser interface {
// 	Parse(string, any) error
// 	MustParse(string, any)
// }

// // TypeParseFns are format functions matched by exact type.
// type TypeParseFns struct {
// 	Int        FormatFn[int]
// 	Int8       FormatFn[int8]
// 	Int16      FormatFn[int16]
// 	Int32      FormatFn[int32] // aka rune
// 	Int64      FormatFn[int64]
// 	Uint       FormatFn[uint]
// 	Uint8      FormatFn[uint8] // aka byte
// 	Uint16     FormatFn[uint16]
// 	Uint32     FormatFn[uint32]
// 	Uint64     FormatFn[uint64]
// 	Uintptr    FormatFn[uintptr]
// 	Float32    FormatFn[float32]
// 	Float64    FormatFn[float64]
// 	Complex64  FormatFn[complex64]
// 	Complex128 FormatFn[complex128]
// 	Bool       FormatFn[bool]
// 	String     FormatFn[string]

// 	// TODO: SliceInt, SliceInt8, ...

// 	// TODO: support custom types
// }

// // KindParseFns are format functions matched by kind (reflect.Kind).
// //
// // Particularly useful for auto-supporting formatting custom types
// // (e.g. String will be applied `type MyString string`
// // if MyString does not have type matched format function)
// type KindParseFns struct {
// 	Int        ParseFn[int]
// 	Int8       ParseFn[int8]
// 	Int16      ParseFn[int16]
// 	Int32      ParseFn[int32] // aka rune
// 	Int64      ParseFn[int64]
// 	Uint       ParseFn[uint]
// 	Uint8      ParseFn[uint8] // aka byte
// 	Uint16     ParseFn[uint16]
// 	Uint32     ParseFn[uint32]
// 	Uint64     ParseFn[uint64]
// 	Uintptr    ParseFn[uintptr]
// 	Float32    ParseFn[float32]
// 	Float64    ParseFn[float64]
// 	Complex64  ParseFn[complex64]
// 	Complex128 ParseFn[complex128]
// 	Bool       ParseFn[bool]
// 	String     ParseFn[string]

// 	Slice  func() ParseFn[reflect.Value]
// 	Struct func() ParseFn[reflect.Value]
// }

// type ParseFns struct {
// 	PerType TypeParseFns
// 	PerKind KindParseFns
// 	// TODO: PerInterface (match when value implements interface)
// }

// type parseConfig struct {
// 	fns             ParseFns
// 	kindPropagation bool
// }

// func NewParser(fns ParseFns, opts ...ParserOption) Parser {
// 	cfg := parseConfig{fns: fns, kindPropagation: false}
// 	for _, opt := range opts {
// 		opt(&cfg)
// 	}
// 	return &parser{cfg: cfg}
// }

// type parser struct {
// 	cfg parseConfig
// }

// // TODO: allow passing reflect.Value here for performance sake?
// func (p *parser) Parse(v string, out any) error {
// 	return p.parse(v, out)
// }

// func (p *parser) MustParse(v string, out any) {
// 	err := p.Parse(v, out)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func (p *parser) parse(v string, out any) error {
// 	str, err, matched := p.parseByType(v)
// 	if matched {
// 		return str, err
// 	}
// 	str, err, matched = p.parseByKind(v)
// 	if matched {
// 		return str, err
// 	}
// 	return "", fmt.Errorf("no parse function for value of type %T", v)
// }
