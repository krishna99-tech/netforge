package cmd

import (
	"fmt"
	"netforge/internal/ping"
	"time"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping [host]",
	Short: "Send ICMP ECHO_REQUEST packets to network hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		count, _ := cmd.Flags().GetInt("count")
		interval, _ := cmd.Flags().GetDuration("interval")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		cfg := ping.Config{
			Count:    count,
			Interval: interval,
			Timeout:  timeout,
		}

		err := ping.Run(host, cfg)
		if err != nil {
			fmt.Printf("Ping failed: %v\n", err)
			fmt.Println("Note: On some systems, running ping requires elevated privileges.")
		}
	},
}

func init() {
	pingCmd.Flags().IntP("count", "c", 4, "Number of ECHO_REQUEST packets to send")
	pingCmd.Flags().DurationP("interval", "i", time.Second, "Wait interval between sending packets")
	pingCmd.Flags().DurationP("timeout", "t", time.Second*5, "Overall timeout for ping operation")
	rootCmd.AddCommand(pingCmd)
}
