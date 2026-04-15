package xstringscase

import (
	"strings"
)

// TODO: add ToFolded for utf8 encoding
// TODO: support configuring how to deal with consequent stop chars (e.g. '__' or '--')
// TODO: support configuring acronyms that should be uppercased (maybe even give control over when to uppercased)
// TODO: support configuring recognizing delimiters (for tokens and statements)
// TODO: support configuring to preserve surrounding delim chars (e.g. spaces)
// TODO: support option for preallocating size for result

// func To(c Case, s string) string {
// 	switch c {
// 	case Snake:
// 		return ToSnake(s)
// 	default:
// 		panic(fmt.Errorf("unexpected case: %s", c.String()))
// 	}
// }

func ToSnake(s string) string {
	tzer := tokenizer{d: []byte(s), i: 0}
	b := strings.Builder{}
	first := true
	for i, j, ok := tzer.Next(); ok; i, j, ok = tzer.Next() {
		if !first {
			b.WriteByte('_')
		}
		b.WriteString(strings.ToLower(s[i:j]))
		first = false
	}
	return b.String()
}
