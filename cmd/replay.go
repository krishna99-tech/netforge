//go:build pcap

package cmd

import (
	"fmt"
	"netforge/internal/replay"

	"github.com/spf13/cobra"
)

var replayCmd = &cobra.Command{
	Use:   "replay [pcap-file]",
	Short: "Replay a PCAP file against a live network interface",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pcapFile := args[0]
		iface, _ := cmd.Flags().GetString("iface")
		filter, _ := cmd.Flags().GetString("filter")
		rateStr, _ := cmd.Flags().GetString("rate")
		loops, _ := cmd.Flags().GetInt("loop")

		rateScale := replay.ParseRateScale(rateStr)

		fmt.Println("Note: Packet replay requires Administrator/root privileges and Npcap/libpcap.")

		opts := replay.Options{
			PcapFile:  pcapFile,
			Interface: iface,
			Filter:    filter,
			RateScale: rateScale,
			LoopCount: loops,
		}

		if err := replay.Run(opts); err != nil {
			fmt.Printf("Replay error: %v\n", err)
		}
	},
}

func init() {
	replayCmd.Flags().StringP("iface", "i", "eth0", "Network interface to replay packets on")
	replayCmd.Flags().StringP("filter", "f", "", "BPF filter to select packets from the PCAP")
	replayCmd.Flags().StringP("rate", "r", "1x", "Replay speed multiplier (e.g., 0.5x, 2x, 10x)")
	replayCmd.Flags().IntP("loop", "l", 1, "Number of times to loop the replay")
	rootCmd.AddCommand(replayCmd)
}
