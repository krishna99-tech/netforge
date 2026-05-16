package mtu

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// DiscoverMTU performs a binary search to find the Path MTU to a target host.
func DiscoverMTU(host string, minMTU, maxMTU int) (int, error) {
	low := minMTU
	high := maxMTU
	best := low

	fmt.Printf("Discovering Path MTU to %s (range: %d - %d)...\n", host, minMTU, maxMTU)

	for low <= high {
		mid := low + (high-low)/2

		// Note: The size we pass to ping is the ICMP payload size.
		// Total IP packet = Payload + ICMP Header (8) + IP Header (20)
		payloadSize := mid - 28

		if payloadSize < 0 {
			break
		}

		success := sendPingProbe(host, payloadSize)
		if success {
			fmt.Printf("[+] MTU %d: Success\n", mid)
			best = mid
			low = mid + 1
		} else {
			fmt.Printf("[-] MTU %d: Fragment needed / Timeout\n", mid)
			high = mid - 1
		}
	}

	return best, nil
}

func sendPingProbe(host string, payloadSize int) bool {
	var cmd *exec.Cmd

	sizeStr := strconv.Itoa(payloadSize)

	if runtime.GOOS == "windows" {
		// -n 1: 1 ping
		// -f: Set Don't Fragment flag in packet
		// -l size: Send buffer size
		// -w 1000: Timeout in ms
		cmd = exec.Command("ping", "-n", "1", "-f", "-l", sizeStr, "-w", "1000", host)
	} else if runtime.GOOS == "darwin" {
		// -c 1: 1 ping
		// -D: Set the Don't Fragment bit
		// -s size: Specify the number of data bytes
		// -W 1000: Time to wait for a response in ms
		cmd = exec.Command("ping", "-c", "1", "-D", "-s", sizeStr, "-W", "1000", host)
	} else {
		// Linux
		// -c 1: 1 ping
		// -M do: Prohibit fragmentation, even local one
		// -s size: Specify the number of data bytes
		// -W 1: Time to wait for a response in seconds
		cmd = exec.Command("ping", "-c", "1", "-M", "do", "-s", sizeStr, "-W", "1", host)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	output := strings.ToLower(string(out))
	if strings.Contains(output, "packet needs to be fragmented") ||
		strings.Contains(output, "frag needed") ||
		strings.Contains(output, "message too long") ||
		strings.Contains(output, "100% packet loss") ||
		strings.Contains(output, "could not fragment") {
		return false
	}

	return true
}
