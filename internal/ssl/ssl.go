package ssl

import (
	"crypto/tls"
	"fmt"
	"time"
)

// CertificateInfo holds details about an SSL certificate.
type CertificateInfo struct {
	Subject    string
	Issuer     string
	Expiry     time.Time
	DNSNames   []string
	IsCritical bool
}

// Inspect retrieves SSL certificate details for a host.
func Inspect(host string, port int) ([]CertificateInfo, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true, // We want to inspect even if invalid
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	var certs []CertificateInfo

	for _, cert := range state.PeerCertificates {
		certs = append(certs, CertificateInfo{
			Subject:  cert.Subject.String(),
			Issuer:   cert.Issuer.String(),
			Expiry:   cert.NotAfter,
			DNSNames: cert.DNSNames,
		})
	}

	return certs, nil
}
