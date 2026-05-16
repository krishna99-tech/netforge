package cmd

import (
	"fmt"
	"netforge/internal/mac"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var macCmd = &cobra.Command{
	Use:   "mac [address]",
	Short: "Lookup the manufacturer (OUI) of a MAC address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		vendor, err := mac.LookupMAC(address)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			res := map[string]string{
				"mac":    address,
				"vendor": vendor,
			}
			utils.PrintJSON(res)
		case "table":
			data := [][]string{
				{"MAC Address", address},
				{"Vendor", vendor},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Printf("\n--- MAC Address Lookup ---\n")
			fmt.Printf("MAC:    %s\n", address)
			fmt.Printf("Vendor: %s\n", vendor)
		}
	},
}

func init() {
	rootCmd.AddCommand(macCmd)
}
