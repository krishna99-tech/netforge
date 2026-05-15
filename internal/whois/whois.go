package whois

import (
	"github.com/likexian/whois"
)

// Query performs a WHOIS lookup for the target domain or IP.
func Query(target string) (string, error) {
	return whois.Whois(target)
}
