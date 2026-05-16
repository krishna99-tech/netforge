package quic

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// ProbeResult holds QUIC/HTTP3 analysis for a host.
type ProbeResult struct {
	Host       string `json:"host"`
	HTTP3      bool   `json:"http3_supported"`
	AltSvc     string `json:"alt_svc_header"`
	StatusCode int    `json:"status_code"`
	Proto      string `json:"protocol_negotiated"`
}

// Probe checks if a host supports HTTP/3 via Alt-Svc header inspection.
// Note: A full QUIC handshake requires quic-go, but Alt-Svc detection via
// standard HTTP/2 is a reliable proxy test.
func Probe(host string, port int) (*ProbeResult, error) {
	url := fmt.Sprintf("https://%s:%d/", host, port)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "NetForge/1.0.1 (QUIC Probe)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer io.ReadAll(resp.Body)
	defer resp.Body.Close()

	altSvc := resp.Header.Get("Alt-Svc")
	http3Supported := false
	if altSvc != "" {
		for _, val := range []string{"h3", "h3-29", "h3-27"} {
			if contains(altSvc, val) {
				http3Supported = true
				break
			}
		}
	}

	return &ProbeResult{
		Host:       fmt.Sprintf("%s:%d", host, port),
		HTTP3:      http3Supported,
		AltSvc:     altSvc,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
	}, nil
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
