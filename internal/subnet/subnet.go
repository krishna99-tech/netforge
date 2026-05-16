package subnet

import (
	"encoding/binary"
	"fmt"
	"net"
)

// SubnetDetails holds the calculated information for a network.
type SubnetDetails struct {
	CIDR           string `json:"cidr"`
	NetworkAddress string `json:"network_address"`
	Broadcast      string `json:"broadcast_address"`
	Netmask        string `json:"netmask"`
	FirstIP        string `json:"first_usable_ip"`
	LastIP         string `json:"last_usable_ip"`
	TotalHosts     uint32 `json:"total_hosts"`
}

// CalculateSubnet computes network details from a CIDR string.
func CalculateSubnet(cidrStr string) (*SubnetDetails, error) {
	ip, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %v", err)
	}

	// Only support IPv4 for now for simplicity
	if ip.To4() == nil {
		return nil, fmt.Errorf("only IPv4 is supported for subnet calculation")
	}

	mask := ipnet.Mask
	network := ipnet.IP.To4()

	// Calculate broadcast
	broadcast := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		broadcast[i] = network[i] | ^mask[i]
	}

	// Calculate host range
	networkUint := binary.BigEndian.Uint32(network)
	broadcastUint := binary.BigEndian.Uint32(broadcast)

	totalHosts := uint32(0)
	firstIP := ""
	lastIP := ""

	if broadcastUint > networkUint+1 {
		totalHosts = (broadcastUint - networkUint) - 1
		firstIP = intToIP(networkUint + 1).String()
		lastIP = intToIP(broadcastUint - 1).String()
	} else {
		// For /31 or /32
		totalHosts = (broadcastUint - networkUint) + 1
		firstIP = network.String()
		lastIP = broadcast.String()
	}

	return &SubnetDetails{
		CIDR:           cidrStr,
		NetworkAddress: network.String(),
		Broadcast:      broadcast.String(),
		Netmask:        net.IP(mask).String(),
		FirstIP:        firstIP,
		LastIP:         lastIP,
		TotalHosts:     totalHosts,
	}, nil
}

func intToIP(i uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, i)
	return ip
}
