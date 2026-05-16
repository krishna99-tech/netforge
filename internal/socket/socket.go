package socket

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// SocketEntry represents a single active socket.
type SocketEntry struct {
	Protocol    string `json:"protocol"`
	LocalAddr   string `json:"local_address"`
	RemoteAddr  string `json:"remote_address"`
	State       string `json:"state"`
	ProcessInfo string `json:"process"`
}

// List returns all active sockets from the OS using native tools.
func List(proto, state string) ([]SocketEntry, error) {
	var out []byte
	var err error

	if runtime.GOOS == "windows" {
		out, err = exec.Command("netstat", "-ano").CombinedOutput()
	} else {
		out, err = exec.Command("ss", "-tulnp").CombinedOutput()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to run socket diagnostics: %v", err)
	}

	return parseOutput(string(out), proto, state, runtime.GOOS), nil
}

func parseOutput(raw, proto, state, os string) []SocketEntry {
	var entries []SocketEntry
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)

		if os == "windows" {
			// netstat -ano output format:
			// Proto  Local Address  Foreign Address  State  PID
			if len(fields) < 4 {
				continue
			}
			p := strings.ToUpper(fields[0])
			if !strings.HasPrefix(p, "TCP") && !strings.HasPrefix(p, "UDP") {
				continue
			}
			if proto != "" && !strings.EqualFold(p, proto) {
				continue
			}
			entry := SocketEntry{
				Protocol:   p,
				LocalAddr:  fields[1],
				RemoteAddr: fields[2],
			}
			if len(fields) >= 5 {
				entry.State = fields[3]
				entry.ProcessInfo = "PID:" + fields[4]
			} else if len(fields) >= 4 {
				entry.ProcessInfo = "PID:" + fields[3]
			}

			if state != "" && !strings.EqualFold(entry.State, state) {
				continue
			}
			entries = append(entries, entry)
		} else {
			// ss -tulnp: Netid State Recv-Q Send-Q Local Peer Process
			if len(fields) < 5 {
				continue
			}
			p := strings.ToLower(fields[0])
			if p == "netid" {
				continue
			}
			if proto != "" && !strings.EqualFold(p, proto) {
				continue
			}
			entry := SocketEntry{
				Protocol:  fields[0],
				State:     fields[1],
				LocalAddr: fields[4],
			}
			if len(fields) > 5 {
				entry.RemoteAddr = fields[5]
			}
			if len(fields) > 6 {
				entry.ProcessInfo = fields[6]
			}

			if state != "" && !strings.EqualFold(entry.State, state) {
				continue
			}
			entries = append(entries, entry)
		}
	}

	return entries
}
