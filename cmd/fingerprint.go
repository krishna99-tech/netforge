//go:build pcap

package cmd

import (
	"fmt"
	"netforge/internal/fingerprint"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var fingerprintCmd = &cobra.Command{
	Use:   "fingerprint [host]",
	Short: "OS fingerprinting via TCP/IP stack analysis",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		port, _ := cmd.Flags().GetInt("port")

		fmt.Printf("Fingerprinting %s (probing port %d)...\n", target, port)
		fmt.Println("Note: Requires Administrator/root and Npcap/libpcap.")

		result, err := fingerprint.Fingerprint(target, port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(result)
		case "table":
			data := [][]string{
				{"IP", result.IP},
				{"TTL", fmt.Sprintf("%d", result.TTL)},
				{"Window Size", fmt.Sprintf("%d", result.WindowSize)},
				{"Guessed OS", result.GuessedOS},
				{"Confidence", result.Confidence},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Println("\n--- OS Fingerprint Result ---")
			fmt.Printf("Target:      %s\n", result.IP)
			fmt.Printf("TTL:         %d\n", result.TTL)
			fmt.Printf("Window Size: %d\n", result.WindowSize)
			fmt.Printf("Guessed OS:  %s\n", result.GuessedOS)
			fmt.Printf("Confidence:  %s\n", result.Confidence)
		}
	},
}

func init() {
	fingerprintCmd.Flags().IntP("port", "p", 80, "Target port to probe (choose an open port for best results)")
	rootCmd.AddCommand(fingerprintCmd)
}
