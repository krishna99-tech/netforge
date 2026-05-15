package cmd

import (
	"fmt"
	"netforge/internal/dns"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns [domain]",
	Short: "Resolve a domain name to IP addresses",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		ips, err := dns.Resolve(domain)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(ips)
		case "table":
			var data [][]string
			for _, ip := range ips {
				data = append(data, []string{ip.String()})
			}
			utils.PrintTable([]string{"IP Address"}, data)
		default:
			fmt.Printf("Results for %s:\n", domain)
			for _, ip := range ips {
				fmt.Printf(" - %s\n", ip)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
