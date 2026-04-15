package xnet

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xptr"
)

func TestLocationString(t *testing.T) {
	testCases := []struct {
		title string
		loc   Location
		url   string
	}{
		{
			"http",
			Location{
				Scheme: "http",
				Host:   "myhost",
				Port:   nil,
			},
			"http://myhost",
		},
		{
			"https",
			Location{
				Scheme: "https",
				Host:   "myhost",
				Port:   nil,
			},
			"https://myhost",
		},
		{
			"subdomain",
			Location{
				Scheme: "https",
				Host:   "sub.myhost",
				Port:   nil,
			},
			"https://sub.myhost",
		},
		{
			"port",
			Location{
				Scheme: "https",
				Host:   "myhost",
				Port:   xptr.T(443),
			},
			"https://myhost:443",
		},
		{
			"port_custom",
			Location{
				Scheme: "https",
				Host:   "myhost",
				Port:   xptr.T(80),
			},
			"https://myhost:80",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			u := tc.loc.URL()
			require.Equal(t, tc.url, u.String())
		})
	}
}
