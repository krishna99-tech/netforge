package traceroute

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// Hop represents a single hop in the traceroute.
type Hop struct {
	HopNumber int    `json:"hop"`
	Address   string `json:"address"`
	Latency   string `json:"latency"`
}

// Trace simulates a traceroute. Due to the complexity of cross-platform raw sockets
// in Go without CGO or requiring root/admin on all platforms, this wraps the system
// traceroute/tracert command and parses the basic output for cross-platform compatibility.
func Trace(host string) ([]Hop, error) {
	// First resolve the host to ensure it exists
	_, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve host: %v", err)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tracert", "-d", "-h", "30", "-w", "1000", host)
	} else {
		cmd = exec.Command("traceroute", "-n", "-m", "30", "-w", "1", host)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("traceroute failed: %v", err)
	}

	return parseTracerouteOutput(string(out), runtime.GOOS), nil
}

func parseTracerouteOutput(output, osType string) []Hop {
	var hops []Hop
	lines := strings.Split(output, "\n")

	// Very simplified parser just to show the structure
	hopCount := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// A very basic heuristic to find lines that look like hops
		if osType == "windows" {
			// Windows tracert output:
			//  1     2 ms     1 ms     1 ms  192.168.1.1
			if strings.Contains(line, " ms ") || strings.Contains(line, "*") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					ipPart := parts[len(parts)-1]
					if ipPart == "timed" || ipPart == "out." {
						ipPart = "* Request timed out."
					}
					hops = append(hops, Hop{
						HopNumber: hopCount,
						Address:   ipPart,
						Latency:   "", // Simplified for now
					})
					hopCount++
				}
			}
		} else {
			// Linux traceroute output
			if strings.Contains(line, " ms") || strings.Contains(line, "*") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					hops = append(hops, Hop{
						HopNumber: hopCount,
						Address:   parts[1],
						Latency:   "",
					})
					hopCount++
				}
			}
		}
	}
	return hops
}
