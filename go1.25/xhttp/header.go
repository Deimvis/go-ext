package xhttp

import (
	"net/http"
)

type ConstHeader interface {
	Get(key string) string
	Values(key string) []string
	// Range stops when f returns false.
	Range(f func(key string, values []string) bool)
	Clone() http.Header
}

func AsConstHeader(h http.Header) ConstHeader {
	return constHeader{Header: h}
}

type constHeader struct {
	http.Header
}

func (ch constHeader) Range(f func(key string, values []string) bool) {
	for key, values := range ch.Header {
		ok := f(key, values)
		if !ok {
			return
		}
	}
}
