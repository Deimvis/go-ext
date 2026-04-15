package xhttp

import (
	"net/http"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

type RoundTripFn func(*http.Request) (*http.Response, error)
type RoundTripWrapFn func(RoundTripFn) RoundTripFn

func WrapRoundTrip(fn RoundTripFn, wraps ...RoundTripWrapFn) RoundTripFn {
	for _, wrap := range wraps {
		fn = wrap(fn)
	}
	return fn
}

func WrapRoundTripper(rt http.RoundTripper, wraps ...RoundTripWrapFn) http.RoundTripper {
	xmust.True(rt != nil)
	rtFn := func(req *http.Request) (*http.Response, error) {
		return rt.RoundTrip(req)
	}
	rtFn = WrapRoundTrip(rtFn, wraps...)

	var underlyingTransport *http.Transport = nil
	if _, ok := rt.(*http.Transport); ok {
		underlyingTransport = rt.(*http.Transport)
	}
	if _, ok := rt.(CustomRoundTripper); ok {
		t, known := rt.(CustomRoundTripper).UnderlyingTransport()
		if known {
			underlyingTransport = t
		}
	}

	return &customRoundTripper{rtFn: rtFn, t: underlyingTransport}
}

type CustomRoundTripper interface {
	http.RoundTripper
	// CustomRoundTripper is decomposed into underlying *http.Transport + wrap functions
	UnderlyingTransport() (*http.Transport, bool)
	Wraps() []RoundTripWrapFn
}

type customRoundTripper struct {
	rtFn RoundTripFn

	t     *http.Transport
	wraps []RoundTripWrapFn
}

func (wrt *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return wrt.rtFn(req)
}

func (wrt *customRoundTripper) UnderlyingTransport() (*http.Transport, bool) {
	if wrt.t != nil {
		return wrt.t, true
	}
	return nil, false
}

func (wrt *customRoundTripper) Wraps() []RoundTripWrapFn {
	return wrt.wraps
}
