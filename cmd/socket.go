package cmd

import (
	"fmt"
	"netforge/internal/socket"
	"netforge/internal/utils"

	"github.com/spf13/cobra"
)

var socketCmd = &cobra.Command{
	Use:   "socket",
	Short: "Socket diagnostics — active connection table",
}

var socketListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active sockets on this machine",
	Run: func(cmd *cobra.Command, args []string) {
		proto, _ := cmd.Flags().GetString("proto")
		state, _ := cmd.Flags().GetString("state")

		entries, err := socket.List(proto, state)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(entries) == 0 {
			fmt.Println("No sockets matched the filter.")
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(entries)
		case "table":
			var data [][]string
			for _, e := range entries {
				data = append(data, []string{e.Protocol, e.LocalAddr, e.RemoteAddr, e.State, e.ProcessInfo})
			}
			utils.PrintTable([]string{"Proto", "Local", "Remote", "State", "Process"}, data)
		default:
			fmt.Println("\n--- Active Sockets ---")
			fmt.Printf("%-6s %-26s %-26s %-14s %s\n", "Proto", "Local", "Remote", "State", "Process")
			fmt.Println(repeat("-", 90))
			for _, e := range entries {
				fmt.Printf("%-6s %-26s %-26s %-14s %s\n", e.Protocol, e.LocalAddr, e.RemoteAddr, e.State, e.ProcessInfo)
			}
		}
	},
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func init() {
	socketListCmd.Flags().StringP("proto", "", "", "Filter by protocol (tcp, udp)")
	socketListCmd.Flags().StringP("state", "", "", "Filter by state (e.g., LISTEN, ESTABLISHED)")
	socketCmd.AddCommand(socketListCmd)
	rootCmd.AddCommand(socketCmd)
}
