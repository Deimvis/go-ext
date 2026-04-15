package xoptional

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xshould"
	"github.com/Deimvis/go-ext/go1.25/xptr"
	yaml_v3 "gopkg.in/yaml.v3"
)

func TestEncodings(t *testing.T) {
	encodings := []encoding{
		{
			"json",
			json.Marshal,
			json.Unmarshal,
			[]byte(`null`),
			nil,
		},
		{
			"yaml_v3",
			yaml_v3.Marshal,
			yaml_v3.Unmarshal,
			[]byte(`null`),
			func(text []byte) []byte {
				if len(text) == 0 || text[len(text)-1] != 0xa {
					text = append(text, 0xa)
				}
				return text
			},
		},
	}
	for _, enc := range encodings {
		runEncodingTestSuite(t, enc)
	}
}

type encoding struct {
	name                  string
	marshal               func(any) ([]byte, error)
	unmarshal             func(data []byte, v any) error
	nullValue             []byte
	marshalPostFormatting func([]byte) []byte
}

func runEncodingTestSuite(
	t *testing.T,
	enc encoding,
) {
	t.Run(enc.name, func(t *testing.T) {
		runMarshalTestSuite(t, enc.marshal, enc)
		runUnmarshalTestSuite(t, enc.unmarshal)
	})
}

func runMarshalTestSuite(
	t *testing.T,
	marshal func(any) ([]byte, error),
	enc encoding,
) {
	t.Run("marshal", func(t *testing.T) {
		for _, tc := range marshalTests {
			t.Run(tc.title, func(t *testing.T) {
				act, err := marshal(tc.v)

				if tc.expErr != nil {
					require.Error(t, err)
					tc.expErr(t, err)
				} else {
					require.NoError(t, err)
					tc.exp(t, act, enc)
				}
			})
		}
	})
}

func runUnmarshalTestSuite(
	t *testing.T,
	unmarshal func([]byte, any) error,
) {
	t.Run("unmarshal", func(t *testing.T) {
		for _, tc := range unmarshalTests {
			t.Run(tc.title, func(t *testing.T) {
				act := tc.v
				err := unmarshal(tc.data, act)
				if tc.expErr != nil {
					require.Error(t, err)
					tc.expErr(t, err)
				} else {
					require.NoError(t, err)
					tc.exp(t, act)
				}
			})
		}
	})
}

var marshalTests = []marshalTc{
	{
		"int/value",
		New[int](42),
		expMarshalValue([]byte(`42`)),
		nil,
	},
	{
		"int/null",
		New[int](),
		expMarshalNull(),
		nil,
	},
	{
		"*int/value",
		New[*int](xptr.T(42)),
		expMarshalValue([]byte(`42`)),
		nil,
	},
	{
		"*int/null",
		New[*int](),
		expMarshalNull(),
		nil,
	},
	// TODO: test for strings quoting
	{
		"string/value",
		New[string]("hello"),
		expMarshalValueByEnc(map[string]marshalExpFn{
			"json":    expMarshalValue([]byte(`"hello"`)),
			"yaml_v3": expMarshalValue([]byte(`hello`)),
		}),
		nil,
	},
	{
		"string/null",
		New[string](),
		expMarshalNull(),
		nil,
	},
	{
		"map/int-values",
		New(map[string]int{
			"key_a": 111,
			"key_b": 222,
		}),
		expMarshalValueByEnc(map[string]marshalExpFn{
			"json": expMarshalValueAny([][]byte{
				[]byte(`{"key_a":111,"key_b":222}`),
				[]byte(`{"key_b":222,"key_a":111}`),
			}),
			"yaml_v3": expMarshalValueAny([][]byte{
				[]byte("key_a: 111\nkey_b: 222"),
				[]byte("key_b: 222\nkey_a: 111"),
			}),
		}),
		nil,
	},
	{
		"map/*int-values/value",
		New(map[string]*int{
			"key_value": xptr.T(123),
		}),
		expMarshalValueByEnc(map[string]marshalExpFn{
			"json":    expMarshalValue([]byte(`{"key_value":123}`)),
			"yaml_v3": expMarshalValue([]byte(`key_value: 123`)),
		}),
		nil,
	},
	{
		"map/*int-values/null",
		New(map[string]*int{
			"key_null": nil,
		}),
		expMarshalValueByEnc(map[string]marshalExpFn{
			"json":    expMarshalValue([]byte(`{"key_null":null}`)),
			"yaml_v3": expMarshalValue([]byte(`key_null: null`)),
		}),
		nil,
	},
	{
		"map/nested",
		New(map[string]map[string]int{
			"l1": {
				"l2": 3,
			},
		}),
		expMarshalValueByEnc(map[string]marshalExpFn{
			"json":    expMarshalValue([]byte(`{"l1":{"l2":3}}`)),
			"yaml_v3": expMarshalValue([]byte("l1:\n    l2: 3")),
		}),
		nil,
	},
}

func expMarshalValue(exp []byte) marshalExpFn {
	return func(t *testing.T, act []byte, enc encoding) {
		if enc.marshalPostFormatting != nil {
			exp = enc.marshalPostFormatting(exp)
		}
		require.Equal(t, exp, act)
	}
}

func expMarshalValueAny(exps [][]byte) marshalExpFn {
	return func(t *testing.T, act []byte, enc encoding) {
		// []byte is not comparable, so workaround with string
		// until we can easily create our own binary predicate
		expsS := ext.Map(exps, func(b []byte) string {
			if enc.marshalPostFormatting != nil {
				b = enc.marshalPostFormatting(b)
			}
			return string(b)
		})
		actS := string(act)
		err := xshould.AnyOf(xcheck.Eq[string](), actS, expsS, xcheck.PrintWhy())
		require.NoError(t, err)
	}
}

func expMarshalValueByEnc(enc2exp map[string]marshalExpFn) marshalExpFn {
	return func(t *testing.T, act []byte, enc encoding) {
		expFn, ok := enc2exp[enc.name]
		xmust.True(ok, "invalid test: not expect value for encoding '%s'", enc.name)
		expFn(t, act, enc)
	}
}

func expMarshalNull() marshalExpFn {
	return func(t *testing.T, act []byte, enc encoding) {
		exp := enc.nullValue
		if enc.marshalPostFormatting != nil {
			exp = enc.marshalPostFormatting(enc.nullValue)
		}
		require.Equal(t, exp, act)
	}
}

var unmarshalTests = []unmarshalTc{
	{
		"int/value",
		[]byte(`42`),
		xptr.T(New[int]()),
		expUnmarshalValue[int](42),
		nil,
	},
	{
		"int/null",
		[]byte(`null`),
		xptr.T(New[int]()),
		expUnmarshalNull[int](),
		nil,
	},
	{
		"*int/value",
		[]byte(`42`),
		xptr.T(New[*int]()),
		expUnmarshalValue[*int](xptr.T(42)),
		nil,
	},
	{
		"*int/null",
		[]byte(`null`),
		xptr.T(New[*int]()),
		expUnmarshalNull[*int](),
		nil,
	},
	{
		"int/error",
		[]byte(`"not number"`),
		xptr.T(New[int]()),
		nil,
		func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
	{
		"string/value",
		[]byte(`"hello"`),
		xptr.T(New[string]()),
		expUnmarshalValue[string](`hello`),
		nil,
	},
	{
		"string/null",
		[]byte(`null`),
		xptr.T(New[string]()),
		expUnmarshalNull[string](),
		nil,
	},
}

func expUnmarshalValue[U any](exp U) func(*testing.T, any) {
	return func(t *testing.T, v any) {
		vv := v.(*T[U])
		require.True(t, vv.HasValue())
		require.Equal(t, exp, vv.Value())
	}
}

func expUnmarshalNull[U any]() func(*testing.T, any) {
	return func(t *testing.T, v any) {
		vv := v.(*T[U])
		require.False(t, vv.HasValue())
	}
}

type marshalTc struct {
	title  string
	v      any
	exp    marshalExpFn
	expErr func(*testing.T, error)
}

type unmarshalTc struct {
	title  string
	data   []byte
	v      any
	exp    unmarshalExpFn
	expErr func(*testing.T, error)
}

type marshalExpFn func(t *testing.T, act []byte, enc encoding)
type unmarshalExpFn func(*testing.T, any)
