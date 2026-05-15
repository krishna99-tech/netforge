package cmd

import (
	"fmt"
	"netforge/internal/tcp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var tcpCmd = &cobra.Command{
	Use:   "tcp [host] [port]",
	Short: "Test TCP connectivity to a host and port",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		portStr := args[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Printf("Invalid port: %s\n", portStr)
			return
		}

		timeout, _ := cmd.Flags().GetDuration("timeout")

		fmt.Printf("Connecting to %s on port %d...\n", host, port)
		result := tcp.Test(host, port, timeout)

		if result.Success {
			fmt.Printf("[+] Connected to %s:%d in %v\n", result.Host, result.Port, result.Latency)
		} else {
			fmt.Printf("[-] Failed to connect to %s:%d: %v\n", result.Host, result.Port, result.Error)
		}
	},
}

func init() {
	tcpCmd.Flags().DurationP("timeout", "t", 5*time.Second, "Timeout for the connection attempt")
	rootCmd.AddCommand(tcpCmd)
}
