package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPublicIP(t *testing.T) {
	TestCases := []struct {
		Name     string
		IP       string
		Expected bool
	}{
		{
			Name:     "Public IPv4",
			IP:       "244.178.44.111",
			Expected: true,
		},
		{
			Name:     "Local IPv4",
			IP:       "127.0.0.1",
			Expected: false,
		},
		{
			Name:     "Private 192*",
			IP:       "192.168.0.1",
			Expected: false,
		},
		{
			Name:     "Private 10*",
			IP:       "10.0.0.1",
			Expected: false,
		},
		{
			Name:     "Private 172*",
			IP:       "172.16.0.0",
			Expected: false,
		},
		{
			Name:     "Public IPv6",
			IP:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			Expected: true,
		},
		{
			Name:     "Private IPv6",
			IP:       "fd00::1",
			Expected: false,
		},
		{
			Name:     "Local IPv6",
			IP:       "fe80::1",
			Expected: false,
		},
		{
			Name:     "Invalid IP",
			IP:       "invalid",
			Expected: false,
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, IsPublicIP(tc.IP))
		})
	}
}
