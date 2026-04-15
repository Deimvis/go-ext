package ext

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHostPort(t *testing.T) {
	testCases := []struct {
		url     string
		expHost string
		expPort *int
	}{
		{
			"http://myhost:8080",
			"myhost",
			ptr(8080),
		},
		{
			"https://myhost:8080",
			"myhost",
			ptr(8080),
		},
		{
			"http://myhost",
			"myhost",
			ptr(80),
		},
		{
			"https://myhost",
			"myhost",
			ptr(443),
		},
		{
			"other://myhost",
			"myhost",
			nil,
		},
	}
	for _, tc := range testCases {
		u, err := url.Parse(tc.url)
		if err != nil {
			panic(err)
		}
		host, port := GetHostPort(u)
		require.Equal(t, tc.expHost, host)
		require.Equal(t, tc.expPort, port)
	}
}
