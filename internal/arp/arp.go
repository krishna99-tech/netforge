//go:build pcap

package arp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	gpcap "github.com/google/gopacket/pcap"
)

// Host represents a discovered network host.
type Host struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

// Scan sends ARP requests to all IPs in the given CIDR and collects replies.
func Scan(iface string, cidr string, timeout time.Duration) ([]Host, error) {
	ips, err := hostsInCIDR(cidr)
	if err != nil {
		return nil, err
	}

	handle, err := gpcap.OpenLive(iface, 65536, true, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface %s: %v", iface, err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("arp"); err != nil {
		return nil, fmt.Errorf("failed to set BPF filter: %v", err)
	}

	netIface, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface: %v", err)
	}

	var addrs []net.Addr
	addrs, err = netIface.Addrs()
	if err != nil || len(addrs) == 0 {
		return nil, fmt.Errorf("no address found on interface %s", iface)
	}

	srcIP := net.ParseIP("0.0.0.0")
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if v4 := ipnet.IP.To4(); v4 != nil {
				srcIP = v4
				break
			}
		}
	}

	// Send ARP requests concurrently
	var wg sync.WaitGroup
	for _, targetIP := range ips {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()
			sendARP(handle, netIface, srcIP, ip)
		}(targetIP)
	}
	wg.Wait()

	// Collect replies
	var hosts []Host
	seen := make(map[string]bool)
	deadline := time.Now().Add(timeout)

	packetSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSrc.Packets() {
		if time.Now().After(deadline) {
			break
		}
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil {
			continue
		}
		arp, _ := arpLayer.(*layers.ARP)
		if arp.Operation == layers.ARPReply {
			ip := net.IP(arp.SourceProtAddress).String()
			mac := net.HardwareAddr(arp.SourceHwAddress).String()
			if !seen[ip] {
				seen[ip] = true
				hosts = append(hosts, Host{IP: ip, MAC: mac})
			}
		}
	}

	return hosts, nil
}

func sendARP(handle *gpcap.Handle, iface *net.Interface, srcIP, dstIP net.IP) {
	eth := &layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(srcIP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(dstIP.To4()),
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err := gopacket.SerializeLayers(buf, opts, eth, arp); err != nil {
		return
	}
	_ = handle.WritePacketData(buf.Bytes())
}

func hostsInCIDR(cidr string) ([]net.IP, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		tmp := make(net.IP, len(ip))
		copy(tmp, ip)
		ips = append(ips, tmp)
	}
	// Trim network and broadcast
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
