package cmd

import (
	"fmt"
	"netforge/internal/netflow"

	"github.com/spf13/cobra"
)

var netflowCmd = &cobra.Command{
	Use:   "netflow",
	Short: "NetFlow v5 flow collector",
}

var netflowListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start a NetFlow v5 listener on a UDP port",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		collector := &netflow.Collector{
			Port: port,
		}

		if err := collector.Listen(); err != nil {
			fmt.Printf("NetFlow error: %v\n", err)
		}
	},
}

func init() {
	netflowListenCmd.Flags().IntP("port", "p", 2055, "UDP port to listen for NetFlow datagrams")
	netflowCmd.AddCommand(netflowListenCmd)
	rootCmd.AddCommand(netflowCmd)
}
