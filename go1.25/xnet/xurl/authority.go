package xurl

import "net/url"

// Authority represents string with <userinfo>@<host>:<port> format.
// Term follows RFC 2396 (https://datatracker.ietf.org/doc/html/rfc2396#section-3.2.2).
type Authority struct {
	// NOTE: I would rather embed url.Userinfo by value
	// in order to avoid indirection for performance sake,
	// but url.Userinfo provides no option to check whether it's empty or not
	// and standard library itself uses it by pointer,
	// so we're forced to do the same.
	*url.Userinfo
	Hostport
}

func (a Authority) String() string {
	if a.Userinfo == nil {
		return a.Hostport.String()
	}
	return a.Userinfo.String() + "@" + a.Hostport.String()
}
