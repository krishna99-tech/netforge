package cmd

import (
	"fmt"
	"netforge/internal/quic"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var quicCmd = &cobra.Command{
	Use:   "quic [host]",
	Short: "Probe a host for QUIC/HTTP3 support via Alt-Svc header",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port, _ := cmd.Flags().GetInt("port")

		fmt.Printf("Probing %s:%d for QUIC/HTTP3 support...\n", host, port)

		result, err := quic.Probe(host, port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		http3Status := "✗ Not supported"
		if result.HTTP3 {
			http3Status = "✓ Supported"
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(result)
		case "table":
			data := [][]string{
				{"Host", result.Host},
				{"HTTP/3 (QUIC)", http3Status},
				{"Alt-Svc Header", result.AltSvc},
				{"Status Code", fmt.Sprintf("%d", result.StatusCode)},
				{"Protocol Used", result.Proto},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Println("\n--- QUIC / HTTP3 Probe ---")
			fmt.Printf("Host:          %s\n", result.Host)
			fmt.Printf("HTTP/3:        %s\n", http3Status)
			fmt.Printf("Alt-Svc:       %s\n", result.AltSvc)
			fmt.Printf("Status Code:   %d\n", result.StatusCode)
			fmt.Printf("Protocol:      %s\n", result.Proto)
		}
	},
}

func init() {
	quicCmd.Flags().IntP("port", "p", 443, "Target port")
	rootCmd.AddCommand(quicCmd)
}
