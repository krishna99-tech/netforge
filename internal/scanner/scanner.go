package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanResult represents the result of a single port scan.
type ScanResult struct {
	Port int
	Open bool
}

// Scanner handles concurrent port scanning.
type Scanner struct {
	Host           string
	Ports          []int
	MaxConcurrency int
	Timeout        time.Duration
}

// Run executes the port scan and returns the open ports.
func (s *Scanner) Run() ([]int, error) {
	fmt.Printf("Scanning %d ports on %s with %d concurrent workers...\n", len(s.Ports), s.Host, s.MaxConcurrency)

	portsToScan := make(chan int, len(s.Ports))
	results := make(chan int)
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
					conn.Close()
					results <- port
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

	foundPorts := []int{}
	for p := range results {
		foundPorts = append(foundPorts, p)
	}

	return foundPorts, nil
}
