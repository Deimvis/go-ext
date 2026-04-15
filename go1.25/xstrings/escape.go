package xstrings

import (
	"fmt"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

func UniqueJoin(sep byte, ss ...string) string {
	return strings.Join(ext.Map(ss, func(s string) string { return Escape(s, sep) }), string(sep))
}

func UniqueSplit(s string, sep byte) []string {
	var ssEsc []string
	begin := 0
	end := 0
	for end < len(s) {
		if s[end] == sep {
			if end+1 < len(s) && s[end+1] == sep {
				end += 2
			} else {
				ssEsc = append(ssEsc, s[begin:end])
				begin = end + 1
				end += 1
			}
		} else {
			end += 1
		}
	}
	ssEsc = append(ssEsc, s[begin:end])
	return ext.Map(ssEsc, func(s string) string { return Unescape(s, sep) })
}

func Escape(s string, escC byte) string {
	b := strings.Builder{}
	for i := range s {
		if s[i] == escC {
			xmust.NoErr(b.WriteByte(escC))
		}
		xmust.NoErr(b.WriteByte(s[i]))
	}
	return b.String()
}

func Unescape(s string, escC byte) string {
	b := strings.Builder{}
	i := 0
	for i < len(s) {
		if s[i] == escC {
			if i+1 < len(s) && s[i+1] == escC {
				xmust.NoErr(b.WriteByte(escC))
				i += 2
			} else {
				panic(fmt.Sprintf("invalid escaping of `%s`, non-escaped char at %d", s, i))
			}
		} else {
			xmust.NoErr(b.WriteByte(s[i]))
			i += 1
		}
	}
	return b.String()
}
