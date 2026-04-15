package xurl

// Host represents either reg-name or ipv4 or ip-literal (ipv6).
// Term follows RFC 2396 (https://datatracker.ietf.org/doc/html/rfc2396#section-3.2.2)
// and changes from RFC 3986 (https://datatracker.ietf.org/doc/html/rfc3986#appendix-D).
type Host struct {
	// TODO: implement proper parsing
	// regname Regname
	// ip      netip.Addr
	v string
}

func (h Host) String() string {
	return h.v
	// if h.ip.IsValid() {
	// 	return h.ip.String()
	// }
	// return h.regname.String()
}
