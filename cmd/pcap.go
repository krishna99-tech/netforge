//go:build pcap

package cmd

import (
	"fmt"
	"netforge/internal/pcap"
	"time"

	"github.com/spf13/cobra"
)

var pcapCmd = &cobra.Command{
	Use:   "pcap",
	Short: "Live packet capture and analysis",
}

var pcapCaptureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture packets on a network interface",
	Run: func(cmd *cobra.Command, args []string) {
		iface, _ := cmd.Flags().GetString("iface")
		filter, _ := cmd.Flags().GetString("filter")
		count, _ := cmd.Flags().GetInt("count")
		durationStr, _ := cmd.Flags().GetString("duration")
		promisc, _ := cmd.Flags().GetBool("promisc")

		duration, _ := time.ParseDuration(durationStr)

		opts := pcap.CaptureOptions{
			Interface:   iface,
			Filter:      filter,
			Count:       count,
			Duration:    duration,
			Promiscuous: promisc,
			SnapLen:     65535,
		}

		err := pcap.RunCapture(opts)
		if err != nil {
			fmt.Printf("Capture failed: %v\n", err)
			fmt.Println("Note: Packet capture requires Administrator/root privileges.")
		}
	},
}

func init() {
	pcapCaptureCmd.Flags().StringP("iface", "i", "eth0", "Network interface to capture on")
	pcapCaptureCmd.Flags().StringP("filter", "f", "", "BPF filter expression")
	pcapCaptureCmd.Flags().IntP("count", "n", 0, "Number of packets to capture (0 = unlimited)")
	pcapCaptureCmd.Flags().String("duration", "0s", "Duration to capture for (e.g., 30s, 1m)")
	pcapCaptureCmd.Flags().Bool("promisc", false, "Enable promiscuous mode")
	
	pcapCmd.AddCommand(pcapCaptureCmd)
	rootCmd.AddCommand(pcapCmd)
}
