package cmd

import (
	"fmt"
	"netforge/internal/traceroute"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var tracerouteCmd = &cobra.Command{
	Use:   "traceroute [host]",
	Short: "Trace the network path to a remote host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		fmt.Printf("Tracing path to %s (this may take a minute)...\n", host)
		
		hops, err := traceroute.Trace(host)
		if err != nil {
			fmt.Printf("Traceroute error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(hops)
		case "table":
			var data [][]string
			for _, hop := range hops {
				data = append(data, []string{
					fmt.Sprintf("%d", hop.HopNumber),
					hop.Address,
				})
			}
			utils.PrintTable([]string{"Hop", "Address"}, data)
		default:
			fmt.Printf("\n--- Traceroute to %s ---\n", host)
			for _, hop := range hops {
				fmt.Printf("%2d  %s\n", hop.HopNumber, hop.Address)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(tracerouteCmd)
}
