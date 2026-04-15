package xstringfb

import "github.com/Deimvis/go-ext/go1.25/xfallback/xfb"

func OnEmpty(s string, fallback string) string {
	empty := func(v string) bool { return len(s) == 0 }
	return xfb.On(empty, s, fallback)
}
