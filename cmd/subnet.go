package cmd

import (
	"fmt"
	"netforge/internal/subnet"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var subnetCmd = &cobra.Command{
	Use:   "subnet [cidr]",
	Short: "Calculate network details from a CIDR (e.g., 192.168.1.0/24)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cidr := args[0]
		details, err := subnet.CalculateSubnet(cidr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(details)
		case "table":
			data := [][]string{
				{"CIDR", details.CIDR},
				{"Network", details.NetworkAddress},
				{"Broadcast", details.Broadcast},
				{"Netmask", details.Netmask},
				{"First Usable", details.FirstIP},
				{"Last Usable", details.LastIP},
				{"Total Hosts", fmt.Sprintf("%d", details.TotalHosts)},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Printf("\n--- Subnet Details for %s ---\n", details.CIDR)
			fmt.Printf("Network:       %s\n", details.NetworkAddress)
			fmt.Printf("Broadcast:     %s\n", details.Broadcast)
			fmt.Printf("Netmask:       %s\n", details.Netmask)
			fmt.Printf("Usable Range:  %s - %s\n", details.FirstIP, details.LastIP)
			fmt.Printf("Total Hosts:   %d\n", details.TotalHosts)
		}
	},
}

func init() {
	rootCmd.AddCommand(subnetCmd)
}
