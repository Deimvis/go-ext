package xurlcfg

import (
	"fmt"
	"net/url"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xshould"
	"github.com/Deimvis/go-ext/go1.25/xoptional"
)

// Origin definition: https://datatracker.ietf.org/doc/html/rfc6454#section-4
type Origin struct {
	Scheme    string `yaml:"scheme" json:"scheme" validate:"required"`
	Authority `yaml:",inline"`
}

type Authority struct {
	Userinfo `yaml:",inline"`
	Hostport `yaml:",inline"`
}

type Userinfo struct {
	Username xoptional.T[string] `json:"username" yaml:"username"`
	Password xoptional.T[string] `json:"password" yaml:"password"`
}

type Hostport struct {
	Host string              `json:"host" yaml:"host" validate:"required"`
	Port xoptional.T[uint16] `json:"port" yaml:"port"`
}

func (o *Origin) URL() *url.URL {
	u := &url.URL{}
	u.Scheme = o.Scheme
	if o.Port.HasValue() {
		u.Host = fmt.Sprintf("%s:%d", o.Host, o.Port.Value())
	} else {
		u.Host = o.Host
	}
	return u
}

func (hp Hostport) ValidateSelf() error {
	if hp.Port.HasValue() {
		err := xshould.Lt(hp.Port.Value(), 65535)
		if err != nil {
			return err
		}
	}
	return nil
}
