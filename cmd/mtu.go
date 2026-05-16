package cmd

import (
	"fmt"
	"netforge/internal/mtu"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var mtuCmd = &cobra.Command{
	Use:   "mtu [host]",
	Short: "Discover Path Maximum Transmission Unit (PMTU)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		min, _ := cmd.Flags().GetInt("min")
		max, _ := cmd.Flags().GetInt("max")

		fmt.Println("Warning: This command sends ICMP packets. It may require Administrator/root privileges depending on your OS.")
		
		bestMtu, err := mtu.DiscoverMTU(host, min, max)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			res := map[string]interface{}{
				"host":     host,
				"path_mtu": bestMtu,
			}
			utils.PrintJSON(res)
		case "table":
			data := [][]string{
				{"Target Host", host},
				{"Path MTU", fmt.Sprintf("%d bytes", bestMtu)},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Printf("\n--- PMTUD Result ---\n")
			fmt.Printf("Target:   %s\n", host)
			fmt.Printf("Path MTU: %d bytes\n", bestMtu)
		}
	},
}

func init() {
	mtuCmd.Flags().IntP("min", "", 576, "Minimum probe size")
	mtuCmd.Flags().IntP("max", "", 9000, "Maximum probe size (up to jumbo frames)")
	rootCmd.AddCommand(mtuCmd)
}
