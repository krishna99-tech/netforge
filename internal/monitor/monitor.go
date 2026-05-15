package monitor

import (
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// GetCPUUsage returns the total CPU usage percentage.
func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		return percentages[0], nil
	}
	return 0, nil
}

// ListeningPort represents a local listening port.
type ListeningPort struct {
	Port uint32
	Type string // TCP or UDP
}

// GetListeningPorts returns a list of local ports currently listening.
func GetListeningPorts() ([]ListeningPort, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	var ports []ListeningPort
	seen := make(map[uint32]bool)

	for _, conn := range conns {
		if conn.Status == "LISTEN" {
			if !seen[conn.Laddr.Port] {
				typeName := "TCP"
				if conn.Type == 2 { // UDP is usually 2, but gopsutil might vary
					typeName = "UDP"
				}
				ports = append(ports, ListeningPort{
					Port: conn.Laddr.Port,
					Type: typeName,
				})
				seen[conn.Laddr.Port] = true
			}
		}
	}

	return ports, nil
}
 
// SystemInfo holds basic system details.
type SystemInfo struct {
	OS            string
	Platform      string
	KernelVersion string
	TotalMemory   uint64
	UsedMemory    uint64
	Uptime        uint64
}
 
// GetSystemInfo returns basic system details.
func GetSystemInfo() (*SystemInfo, error) {
	h, err := host.Info()
	if err != nil {
		return nil, err
	}
 
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
 
	return &SystemInfo{
		OS:            h.OS,
		Platform:      h.Platform,
		KernelVersion: h.KernelVersion,
		TotalMemory:   v.Total,
		UsedMemory:    v.Used,
		Uptime:        h.Uptime,
	}, nil
}
 
// BandwidthInfo holds upload and download rates.
type BandwidthInfo struct {
	SentRate float64 // bytes/sec
	RecvRate float64 // bytes/sec
}
 
// GetBandwidthUsage measures network bandwidth over a 1-second interval.
func GetBandwidthUsage() (*BandwidthInfo, error) {
	io1, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
 
	time.Sleep(time.Second)
 
	io2, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
 
	if len(io1) > 0 && len(io2) > 0 {
		return &BandwidthInfo{
			SentRate: float64(io2[0].BytesSent - io1[0].BytesSent),
			RecvRate: float64(io2[0].BytesRecv - io1[0].BytesRecv),
		}, nil
	}
 
	return &BandwidthInfo{}, nil
}

// InterfaceInfo holds network interface details.
type InterfaceInfo struct {
	Name        string
	IPAddresses []string
	MACAddress  string
	Flags       []string
}

// GetNetworkInterfaces returns details for all active network interfaces.
func GetNetworkInterfaces() ([]InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var infos []InterfaceInfo
	for _, iface := range interfaces {
		var ips []string
		for _, addr := range iface.Addrs {
			ips = append(ips, addr.Addr)
		}

		infos = append(infos, InterfaceInfo{
			Name:        iface.Name,
			IPAddresses: ips,
			MACAddress:  iface.HardwareAddr,
			Flags:       iface.Flags,
		})
	}
	return infos, nil
}

// GetWiFiSSID attempts to retrieve the active WiFi SSID (Windows only via netsh).
func GetWiFiSSID() string {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "N/A (Non-Windows or No WiFi)"
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for a line that contains "SSID" but is not the "BSSID" line
		if strings.Contains(line, "SSID") && !strings.Contains(line, "BSSID") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return "Disconnected"
}
