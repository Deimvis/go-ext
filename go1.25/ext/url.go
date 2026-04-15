package ext

import (
	"net/url"
	"strconv"
	"strings"
)

// GetHostPort unlike u.Port() returns integer
// and defaults to 80 or 443 port
// according to url scheme.
// TODO: deprecated; do not use this function
// TODO: make and use function to infer scheme from port instead
func GetHostPort(u *url.URL) (string, *int) {
	u.Hostname()
	u.Port()
	var host string
	var port *int
	if strings.Contains(u.Host, ":") {
		parts := strings.Split(u.Host, ":")
		if len(parts) != 2 {
			panic("url host is invalid")
		}
		host = parts[0]
		portValue, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			panic(err)
		}
		tmp := int(portValue)
		port = &tmp
	} else {
		host = u.Host
		switch u.Scheme {
		case "http":
			tmp := 80
			port = &tmp
		case "https":
			tmp := 443
			port = &tmp
		default:
			port = nil
		}
	}
	return host, port
}
