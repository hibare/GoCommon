// Package ip provides utilities for IP address operations.
package ip

import "net"

// IsPublicIP checks if the IP address is public.
func IsPublicIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check if the IP address is not a loopback or link-local address
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return false
	}

	return true
}
