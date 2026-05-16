//go:build pcap

package fingerprint

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	gpcap "github.com/google/gopacket/pcap"
)

// Result holds the OS fingerprinting analysis.
type Result struct {
	IP              string `json:"ip"`
	TTL             uint8  `json:"ttl"`
	WindowSize      uint16 `json:"window_size"`
	GuessedOS       string `json:"guessed_os"`
	Confidence      string `json:"confidence"`
}

// Fingerprint sends a TCP SYN to target:port and analyses the SYN-ACK response.
func Fingerprint(target string, port int) (*Result, error) {
	iface, srcIP, err := getOutboundInterface(target)
	if err != nil {
		return nil, fmt.Errorf("could not determine outbound interface: %v", err)
	}

	handle, err := gpcap.OpenLive(iface, 65536, false, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %v", err)
	}
	defer handle.Close()

	filter := fmt.Sprintf("tcp and src host %s and src port %d", target, port)
	if err := handle.SetBPFFilter(filter); err != nil {
		return nil, fmt.Errorf("failed to set BPF filter: %v", err)
	}

	// Send a SYN packet
	if err := sendSYN(handle, srcIP, target, port); err != nil {
		return nil, fmt.Errorf("failed to send SYN: %v", err)
	}

	// Listen for SYN-ACK
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	deadline := time.Now().Add(3 * time.Second)

	for packet := range src.Packets() {
		if time.Now().After(deadline) {
			break
		}
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if ipLayer == nil || tcpLayer == nil {
			continue
		}
		ip, _ := ipLayer.(*layers.IPv4)
		tcp, _ := tcpLayer.(*layers.TCP)

		if tcp.SYN && tcp.ACK {
			return analyzeResponse(target, ip.TTL, tcp.Window), nil
		}
	}

	return nil, fmt.Errorf("no SYN-ACK received from %s:%d (host may be down or port closed)", target, port)
}

func analyzeResponse(ip string, ttl uint8, window uint16) *Result {
	guessedOS := "Unknown"
	confidence := "Low"

	switch {
	case ttl >= 120 && ttl <= 128:
		guessedOS = "Windows"
		confidence = "High"
	case ttl >= 60 && ttl <= 64:
		guessedOS = "Linux / Android"
		confidence = "High"
	case ttl >= 250 && ttl <= 255:
		guessedOS = "Cisco IOS / Network Device"
		confidence = "High"
	case ttl >= 120 && ttl < 240:
		guessedOS = "macOS / BSD"
		confidence = "Medium"
	}

	// Refine with window size
	switch {
	case window == 65535:
		guessedOS = "macOS / BSD"
		confidence = "High"
	case window == 8192:
		if guessedOS != "Windows" {
			guessedOS = "Windows (older)"
		}
		confidence = "Medium"
	case window == 29200 || window == 14600:
		guessedOS = "Linux"
		confidence = "High"
	}

	return &Result{
		IP:         ip,
		TTL:        ttl,
		WindowSize: window,
		GuessedOS:  guessedOS,
		Confidence: confidence,
	}
}

func sendSYN(handle *gpcap.Handle, srcIP, dstIP string, dstPort int) error {
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0x00, 0x01},
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		SrcIP:    net.ParseIP(srcIP),
		DstIP:    net.ParseIP(dstIP),
		Protocol: layers.IPProtocolTCP,
		TTL:      64,
		Version:  4,
		IHL:      5,
	}
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(50000),
		DstPort: layers.TCPPort(dstPort),
		SYN:     true,
		Seq:     12345678,
		Window:  65535,
		Options: []layers.TCPOption{
			{OptionType: layers.TCPOptionKindMSS, OptionLength: 4, OptionData: []byte{0x05, 0xb4}},
			{OptionType: layers.TCPOptionKindSACKPermitted, OptionLength: 2},
			{OptionType: layers.TCPOptionKindWindowScale, OptionLength: 3, OptionData: []byte{0x06}},
		},
	}
	_ = tcp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err := gopacket.SerializeLayers(buf, opts, eth, ip, tcp); err != nil {
		return err
	}
	return handle.WritePacketData(buf.Bytes())
}

func getOutboundInterface(target string) (string, string, error) {
	conn, err := net.Dial("udp", target+":80")
	if err != nil {
		return "", "", err
	}
	defer conn.Close()
	srcIP := conn.LocalAddr().(*net.UDPAddr).IP.String()

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.String() == srcIP {
				return i.Name, srcIP, nil
			}
		}
	}
	return "eth0", srcIP, nil
}
