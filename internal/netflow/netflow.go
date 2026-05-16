package netflow

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// FlowRecord represents a decoded NetFlow v5 flow record.
type FlowRecord struct {
	SrcIP    string `json:"src_ip"`
	DstIP    string `json:"dst_ip"`
	SrcPort  uint16 `json:"src_port"`
	DstPort  uint16 `json:"dst_port"`
	Protocol uint8  `json:"protocol"`
	Packets  uint32 `json:"packets"`
	Bytes    uint32 `json:"bytes"`
}

// Collector listens for NetFlow datagrams on a UDP port.
type Collector struct {
	Port      int
	Handler   func(FlowRecord)
}

// Listen starts a UDP listener and decodes incoming NetFlow v5 packets.
func (c *Collector) Listen() error {
	addr := fmt.Sprintf("0.0.0.0:%d", c.Port)
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP %s: %v", addr, err)
	}
	defer pc.Close()

	fmt.Printf("NetFlow collector listening on UDP port %d...\n", c.Port)
	fmt.Println("Press Ctrl+C to stop.")

	buf := make([]byte, 65535)
	for {
		n, remote, err := pc.ReadFrom(buf)
		if err != nil {
			return err
		}

		records, err := decodeV5(buf[:n])
		if err != nil {
			fmt.Printf("[WARN] Failed to decode packet from %s: %v\n", remote, err)
			continue
		}

		for _, r := range records {
			if c.Handler != nil {
				c.Handler(r)
			} else {
				printRecord(r)
			}
		}
	}
}

// decodeV5 decodes a NetFlow v5 UDP datagram into flow records.
func decodeV5(data []byte) ([]FlowRecord, error) {
	// NetFlow v5 header is 24 bytes
	if len(data) < 24 {
		return nil, fmt.Errorf("packet too short for NetFlow v5 header")
	}

	version := binary.BigEndian.Uint16(data[0:2])
	if version != 5 {
		return nil, fmt.Errorf("unsupported NetFlow version: %d (only v5 supported)", version)
	}

	count := binary.BigEndian.Uint16(data[2:4])
	var records []FlowRecord

	// Each v5 record is 48 bytes, starting at offset 24
	for i := uint16(0); i < count; i++ {
		offset := 24 + int(i)*48
		if offset+48 > len(data) {
			break
		}
		r := data[offset : offset+48]
		records = append(records, FlowRecord{
			SrcIP:    net.IP(r[0:4]).String(),
			DstIP:    net.IP(r[4:8]).String(),
			SrcPort:  binary.BigEndian.Uint16(r[32:34]),
			DstPort:  binary.BigEndian.Uint16(r[34:36]),
			Protocol: r[38],
			Packets:  binary.BigEndian.Uint32(r[16:20]),
			Bytes:    binary.BigEndian.Uint32(r[20:24]),
		})
	}
	return records, nil
}

func printRecord(r FlowRecord) {
	proto := protoName(r.Protocol)
	fmt.Printf("[%s] %s:%d -> %s:%d  pkts=%d bytes=%d\n",
		time.Now().Format("15:04:05"), r.SrcIP, r.SrcPort, r.DstIP, r.DstPort, r.Packets, r.Bytes)
	_ = proto
}

func protoName(proto uint8) string {
	switch proto {
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	case 1:
		return "ICMP"
	default:
		return fmt.Sprintf("proto/%d", proto)
	}
}
