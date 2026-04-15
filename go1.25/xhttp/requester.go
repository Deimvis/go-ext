package xhttp

import (
	"net/http"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"go.uber.org/fx"
)

// https://github.com/golang/go/issues/16047
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
	GetTransport() http.RoundTripper
	Clone() RequesterReflect
}

type RequesterReflect interface {
	Requester
	SetTransport(http.RoundTripper)
}

type RequesterParams struct {
	fx.In

	Client       *http.Client
	GlobalParams *RequesterGlobalParams `optional:"true"`
}

type RequesterGlobalParams struct {
	// RoundTripWraps are useful to implement request hooks (e.g. logging or metrics collection)
	// Some examples of [RoundTripWrapFn] implementations:
	// - Logging with Zap
	// - Metrics collection with Prometheus
	RoundTripWraps []RoundTripWrapFn `optional:"true"`
}

func NewRequester(p RequesterParams) Requester {
	xmust.NotNilPtr(p.Client, "nil *http.Client")
	client := p.Client
	if p.GlobalParams != nil {
		gp := p.GlobalParams
		if client.Transport == nil {
			client.Transport = http.DefaultTransport
		}
		client.Transport = WrapRoundTripper(client.Transport, gp.RoundTripWraps...)
	}

	return &requester{client: client}
}

type requester struct {
	client *http.Client
}

func (r *requester) Do(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
}

func (r *requester) GetTransport() http.RoundTripper {
	return r.client.Transport
}

func (r *requester) Clone() RequesterReflect {
	panic("TODO: implement http.Client clone")
}

func (r *requester) SetTransport(t http.RoundTripper) {
	r.client.Transport = t
}
