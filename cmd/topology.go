//go:build pcap

package cmd

import (
	"fmt"
	"netforge/internal/arp"
	"netforge/internal/utils"
	"time"

	"github.com/spf13/cobra"
)

var topologyCmd = &cobra.Command{
	Use:   "topology",
	Short: "Network topology and host discovery (ARP)",
}

var topologyScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "ARP scan to discover all live hosts on a subnet",
	Run: func(cmd *cobra.Command, args []string) {
		iface, _ := cmd.Flags().GetString("iface")
		cidr, _ := cmd.Flags().GetString("cidr")
		timeoutSec, _ := cmd.Flags().GetInt("timeout")

		if cidr == "" {
			fmt.Println("Error: --cidr is required (e.g., --cidr 192.168.1.0/24)")
			return
		}

		fmt.Printf("ARP scanning %s on interface %s...\n", cidr, iface)
		hosts, err := arp.Scan(iface, cidr, time.Duration(timeoutSec)*time.Second)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Note: This command requires Administrator/root privileges and Npcap/libpcap.")
			return
		}

		if len(hosts) == 0 {
			fmt.Println("No live hosts discovered.")
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(hosts)
		case "table":
			var data [][]string
			for _, h := range hosts {
				data = append(data, []string{h.IP, h.MAC})
			}
			utils.PrintTable([]string{"IP Address", "MAC Address"}, data)
		default:
			fmt.Printf("\n--- Discovered Hosts on %s ---\n", cidr)
			for _, h := range hosts {
				fmt.Printf("[+] %-18s  %s\n", h.IP, h.MAC)
			}
		}
	},
}

func init() {
	topologyScanCmd.Flags().StringP("iface", "i", "eth0", "Network interface to scan from")
	topologyScanCmd.Flags().StringP("cidr", "", "", "Target subnet in CIDR notation (e.g., 192.168.1.0/24)")
	topologyScanCmd.Flags().IntP("timeout", "t", 2, "Seconds to wait for ARP replies")
	topologyCmd.AddCommand(topologyScanCmd)
	rootCmd.AddCommand(topologyCmd)
}
