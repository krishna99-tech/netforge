package cmd

import (
	"fmt"
	"netforge/internal/scanner"
	"time"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [host]",
	Short: "Check for open common ports on a host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
		
		// Common ports
		ports := []int{21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5900, 8000, 8080, 8443}

		s := scanner.Scanner{
			Host:           host,
			Ports:          ports,
			MaxConcurrency: maxConcurrency,
			Timeout:        500 * time.Millisecond,
		}

		foundPorts, err := s.Run()
		if err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			return
		}

		if len(foundPorts) == 0 {
			fmt.Println("No open ports found among the common ones.")
		} else {
			fmt.Println("\n--- Open Ports ---")
			for _, p := range foundPorts {
				fmt.Printf("[+] Port %d is OPEN\n", p)
			}
		}
	},
}

func init() {
	scanCmd.Flags().IntP("max-concurrency", "c", 100, "Maximum number of concurrent port scans")
	rootCmd.AddCommand(scanCmd)
}
