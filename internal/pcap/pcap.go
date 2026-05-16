//go:build pcap

package pcap

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	gpcap "github.com/google/gopacket/pcap"
)

// CaptureOptions defines the settings for packet capture
type CaptureOptions struct {
	Interface   string
	Filter      string
	Count       int
	Duration    time.Duration
	Promiscuous bool
	SnapLen     int
}

// RunCapture starts a live packet capture on the specified interface
func RunCapture(opts CaptureOptions) error {
	fmt.Printf("Starting capture on %s...\n", opts.Interface)

	handle, err := gpcap.OpenLive(opts.Interface, int32(opts.SnapLen), opts.Promiscuous, gpcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open interface: %v", err)
	}
	defer handle.Close()

	if opts.Filter != "" {
		if err := handle.SetBPFFilter(opts.Filter); err != nil {
			return fmt.Errorf("failed to set BPF filter: %v", err)
		}
		fmt.Printf("BPF Filter applied: %s\n", opts.Filter)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	
	packetCount := 0
	startTime := time.Now()

	for packet := range packetSource.Packets() {
		decodePacket(packet)
		packetCount++

		if opts.Count > 0 && packetCount >= opts.Count {
			fmt.Printf("Reached packet limit (%d)\n", opts.Count)
			break
		}

		if opts.Duration > 0 && time.Since(startTime) >= opts.Duration {
			fmt.Printf("Reached time limit (%v)\n", opts.Duration)
			break
		}
	}

	fmt.Printf("\nCapture complete. Captured %d packets.\n", packetCount)
	return nil
}

func decodePacket(packet gopacket.Packet) {
	fmt.Printf("\n--- Packet [%s] ---\n", packet.Metadata().Timestamp.Format(time.RFC3339))
	
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		fmt.Printf("IPv4: %s -> %s\n", ip.SrcIP, ip.DstIP)
	} else if ipLayer := packet.Layer(layers.LayerTypeIPv6); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv6)
		fmt.Printf("IPv6: %s -> %s\n", ip.SrcIP, ip.DstIP)
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Printf("TCP: %d -> %d (Seq: %d)\n", tcp.SrcPort, tcp.DstPort, tcp.Seq)
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		fmt.Printf("UDP: %d -> %d\n", udp.SrcPort, udp.DstPort)
	} else if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		icmp, _ := icmpLayer.(*layers.ICMPv4)
		fmt.Printf("ICMP: Type %d Code %d\n", icmp.TypeCode.Type(), icmp.TypeCode.Code())
	}
}
