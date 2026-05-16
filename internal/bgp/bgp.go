package bgp

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type ASNInfo struct {
	IP        string `json:"ip"`
	ASN       string `json:"asn"`
	Prefix    string `json:"bgp_prefix"`
	Country   string `json:"country"`
	Registry  string `json:"registry"`
	Allocated string `json:"allocated"`
}

// reverseIP reverses an IPv4 address for DNS querying.
func reverseIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s.%s", parts[3], parts[2], parts[1], parts[0])
}

// LookupASN queries origin.asn.cymru.com for BGP info.
func LookupASN(ip string) (*ASNInfo, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To4() == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	reversed := reverseIP(ip)
	query := reversed + ".origin.asn.cymru.com."

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(query, dns.TypeTXT)

	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %v", err)
	}

	if len(r.Answer) == 0 {
		return nil, fmt.Errorf("no BGP information found for IP")
	}

	txtRecord, ok := r.Answer[0].(*dns.TXT)
	if !ok || len(txtRecord.Txt) == 0 {
		return nil, fmt.Errorf("invalid TXT record format")
	}

	// Format: "15169 | 8.8.8.0/24 | US | arin | 1992-12-01"
	data := txtRecord.Txt[0]
	parts := strings.Split(data, " | ")

	if len(parts) < 5 {
		return nil, fmt.Errorf("unexpected record format: %s", data)
	}

	return &ASNInfo{
		IP:        ip,
		ASN:       "AS" + strings.TrimSpace(parts[0]),
		Prefix:    strings.TrimSpace(parts[1]),
		Country:   strings.TrimSpace(parts[2]),
		Registry:  strings.TrimSpace(parts[3]),
		Allocated: strings.TrimSpace(parts[4]),
	}, nil
}
