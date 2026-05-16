package cmd

import (
	"fmt"
	"netforge/internal/ip"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP address and Geolocation utilities",
}

var ipInfoCmd = &cobra.Command{
	Use:   "info [ip/domain]",
	Short: "Get geolocation and ISP info for an IP or domain",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := ""
		if len(args) > 0 {
			target = args[0]
		}

		info, err := ip.GetIPInfo(target)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(info)
		case "table":
			data := [][]string{
				{"IP", info.Query},
				{"Country", info.Country},
				{"City", info.City},
				{"ISP", info.ISP},
				{"ASN", info.AS},
				{"Timezone", info.Timezone},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Printf("\n--- IP Information [%s] ---\n", info.Query)
			fmt.Printf("Location:  %s, %s, %s\n", info.City, info.Region, info.Country)
			fmt.Printf("ISP:       %s\n", info.ISP)
			fmt.Printf("ASN:       %s\n", info.AS)
			fmt.Printf("Timezone:  %s\n", info.Timezone)
			fmt.Printf("Coords:    %f, %f\n", info.Lat, info.Lon)
		}
	},
}

var ipPublicCmd = &cobra.Command{
	Use:   "public",
	Short: "Show your current public IP address",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := ip.GetIPInfo("")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Your Public IP: %s (%s, %s)\n", info.Query, info.City, info.Country)
	},
}

func init() {
	ipCmd.AddCommand(ipInfoCmd)
	ipCmd.AddCommand(ipPublicCmd)
	rootCmd.AddCommand(ipCmd)
}
