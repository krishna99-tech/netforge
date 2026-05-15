package dns

import (
	"net"
)

// Resolve returns the IP addresses for a given domain.
func Resolve(domain string) ([]net.IP, error) {
	return net.LookupIP(domain)
}
