package cmd

import (
	"fmt"
	"netforge/internal/bgp"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var bgpCmd = &cobra.Command{
	Use:   "bgp",
	Short: "BGP and Autonomous System Number (ASN) lookups",
}

var bgpAsnCmd = &cobra.Command{
	Use:   "asn [ip]",
	Short: "Resolve ASN and BGP prefix for an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		info, err := bgp.LookupASN(ip)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(info)
		case "table":
			data := [][]string{
				{"IP", info.IP},
				{"ASN", info.ASN},
				{"BGP Prefix", info.Prefix},
				{"Country", info.Country},
				{"Registry", info.Registry},
				{"Allocated", info.Allocated},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Println("\n--- BGP / ASN Information ---")
			fmt.Printf("IP:          %s\n", info.IP)
			fmt.Printf("ASN:         %s\n", info.ASN)
			fmt.Printf("BGP Prefix:  %s\n", info.Prefix)
			fmt.Printf("Country:     %s\n", info.Country)
			fmt.Printf("Registry:    %s\n", info.Registry)
			fmt.Printf("Allocated:   %s\n", info.Allocated)
		}
	},
}

func init() {
	bgpCmd.AddCommand(bgpAsnCmd)
	rootCmd.AddCommand(bgpCmd)
}
