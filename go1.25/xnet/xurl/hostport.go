package xurl

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/xnet/xurl/internal/stdurl"
)

// Hostport represents string with <host>:<port> format.
// Term follows RFC 2396 (https://datatracker.ietf.org/doc/html/rfc2396#section-3.2.2).
// It's an alternative for url.URL.Host field.
type Hostport struct {
	Host Host
	Port string
}

// PortInt returns port as integer.
// PortInt returns nil when no port
// and returns false if port is not integer
// (when no port it returns true).
// TODO: use uint16 for port
func (hp Hostport) PortInt() (*int, bool) {
	if len(hp.Port) == 0 {
		return nil, true
	}
	p, err := strconv.ParseInt(hp.Port, 10, 64)
	if err != nil {
		return nil, false
	}
	pi := int(p)
	return &pi, true
}

func (hp Hostport) String() string {
	if hp.Port == "" {
		return hp.Host.String()
	}
	return hp.Host.String() + ":" + hp.Port
}

func ParseHostport(s string) (Hostport, error) {
	v, err := stdurl.ParseHost(s)
	if err != nil {
		return Hostport{}, err
	}
	var host string = v
	var port string
	if len(v) > 0 && v[0] == '[' {
		i := strings.Index(v, "]")
		if i < 0 {
			return Hostport{}, errors.New("invalid hostport: no ']' symbol")
		}
		if i+1 < len(v) {
			if v[i+1] != ':' {
				return Hostport{}, errors.New("invalid hostport: no ':' after ']'")
			}
			host = v[:i+1]
			port = v[i+2:]
		}
	} else if len(v) > 0 {
		i := strings.Index(v, ":")
		if i >= 0 {
			host = v[:i]
			port = v[i+1:]
		}
	}
	return Hostport{Host: Host{v: host}, Port: port}, err
}
