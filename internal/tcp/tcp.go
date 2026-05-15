package tcp

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a TCP connection test.
type Result struct {
	Host    string
	Port    int
	Success bool
	Latency time.Duration
	Error   error
}

// Test checks connectivity to a host on a specific port.
func Test(host string, port int, timeout time.Duration) Result {
	address := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Host:    host,
			Port:    port,
			Success: false,
			Latency: latency,
			Error:   err,
		}
	}
	conn.Close()

	return Result{
		Host:    host,
		Port:    port,
		Success: true,
		Latency: latency,
	}
}
