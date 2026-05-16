package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// ScanResult represents the result of a single port scan.
type ScanResult struct {
	Port   int    `json:"port"`
	Open   bool   `json:"open"`
	Banner string `json:"banner,omitempty"`
}

// Scanner handles concurrent port scanning.
type Scanner struct {
	Host           string
	Ports          []int
	MaxConcurrency int
	Timeout        time.Duration
	GrabBanner     bool
}

// Run executes the port scan and returns the open ports.
func (s *Scanner) Run() ([]ScanResult, error) {
	fmt.Printf("Scanning %d ports on %s with %d concurrent workers...\n", len(s.Ports), s.Host, s.MaxConcurrency)

	portsToScan := make(chan int, len(s.Ports))
	results := make(chan ScanResult)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < s.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portsToScan {
				address := fmt.Sprintf("%s:%d", s.Host, port)
				conn, err := net.DialTimeout("tcp", address, s.Timeout)
				if err == nil {
					// Port is open
					banner := ""
					if s.GrabBanner {
						// Attempt to read a banner
						_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
						buf := make([]byte, 1024)
						n, readErr := conn.Read(buf)
						if readErr == nil && n > 0 {
							banner = strings.TrimSpace(string(buf[:n]))
						}
					}
					conn.Close()
					results <- ScanResult{Port: port, Open: true, Banner: banner}
				}
			}
		}()
	}

	// Populate the portsToScan channel
	for _, port := range s.Ports {
		portsToScan <- port
	}
	close(portsToScan)

	// Goroutine to wait for all workers and then close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	foundPorts := []ScanResult{}
	for r := range results {
		foundPorts = append(foundPorts, r)
	}

	return foundPorts, nil
}
