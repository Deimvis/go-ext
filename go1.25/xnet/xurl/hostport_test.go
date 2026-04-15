package xurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseHostport(t *testing.T) {
	tcs := []struct {
		title   string
		s       string
		expHost string
		expPort string
		expErr  *string
	}{
		{
			"host-port",
			"myhost:123",
			"myhost",
			"123",
			nil,
		},
		{
			"host-only",
			"myhost",
			"myhost",
			"",
			nil,
		},
		{
			"port-only",
			":123",
			"",
			"123",
			nil,
		},
		{
			"ipv4-port",
			"1.2.3.4:123",
			"1.2.3.4",
			"123",
			nil,
		},
		{
			"ipv6-port",
			"[1234:::1234:1234:::1234]:123",
			"[1234:::1234:1234:::1234]",
			"123",
			nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			hp, err := ParseHostport(tc.s)
			if tc.expErr != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), *tc.expErr)
			} else {
				require.Equal(t, tc.expHost, hp.Host.String())
				require.Equal(t, tc.expPort, hp.Port)
			}
		})
	}
}
