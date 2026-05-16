package tls

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// CipherInfo holds details about a single TLS cipher/version combination.
type CipherInfo struct {
	Version     string `json:"version"`
	Supported   bool   `json:"supported"`
	CipherSuite string `json:"cipher_suite,omitempty"`
}

// AuditResult holds the complete TLS audit for a host.
type AuditResult struct {
	Host        string       `json:"host"`
	Grade       string       `json:"grade"`
	Certificate CertInfo     `json:"certificate"`
	Versions    []CipherInfo `json:"tls_versions"`
}

// CertInfo holds certificate details.
type CertInfo struct {
	Subject   string    `json:"subject"`
	Issuer    string    `json:"issuer"`
	ExpiresAt time.Time `json:"expires_at"`
	DaysLeft  int       `json:"days_remaining"`
}

var tlsVersions = []struct {
	Name    string
	Version uint16
}{
	{"TLS 1.3", tls.VersionTLS13},
	{"TLS 1.2", tls.VersionTLS12},
	{"TLS 1.1", tls.VersionTLS11},
	{"TLS 1.0", tls.VersionTLS10},
}

// Audit performs TLS version and certificate inspection against a host.
func Audit(host string, port int) (*AuditResult, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	result := &AuditResult{Host: address}

	// Test each TLS version
	for _, v := range tlsVersions {
		info := CipherInfo{Version: v.Name}

		conf := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         v.Version,
			MaxVersion:         v.Version,
			ServerName:         host,
		}
		conn, err := tls.DialWithDialer(
			&net.Dialer{Timeout: 4 * time.Second},
			"tcp", address, conf,
		)
		if err == nil {
			info.Supported = true
			info.CipherSuite = tls.CipherSuiteName(conn.ConnectionState().CipherSuite)
			conn.Close()
		}
		result.Versions = append(result.Versions, info)
	}

	// Get certificate info using a fresh connection
	conf := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 4 * time.Second}, "tcp", address, conf)
	if err == nil {
		defer conn.Close()
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) > 0 {
			leaf := certs[0]
			result.Certificate = CertInfo{
				Subject:   leaf.Subject.CommonName,
				Issuer:    leaf.Issuer.CommonName,
				ExpiresAt: leaf.NotAfter,
				DaysLeft:  int(time.Until(leaf.NotAfter).Hours() / 24),
			}
		}
	}

	result.Grade = computeGrade(result)
	return result, nil
}

// computeGrade assigns a simplified security grade based on TLS support.
func computeGrade(r *AuditResult) string {
	tls13Supported := false
	tls12Supported := false
	tls11OrLower := false

	for _, v := range r.Versions {
		if v.Version == "TLS 1.3" && v.Supported {
			tls13Supported = true
		}
		if v.Version == "TLS 1.2" && v.Supported {
			tls12Supported = true
		}
		if (v.Version == "TLS 1.1" || v.Version == "TLS 1.0") && v.Supported {
			tls11OrLower = true
		}
	}

	if tls13Supported && !tls11OrLower {
		return "A+"
	}
	if tls12Supported && !tls11OrLower {
		return "A"
	}
	if tls12Supported && tls11OrLower {
		return "B"
	}
	if tls11OrLower {
		return "C"
	}
	return "F"
}
