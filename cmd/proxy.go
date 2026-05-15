package cmd

import (
	"fmt"
	"netforge/internal/proxy"
	"os"

	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Network proxy utilities",
}

var proxyStartCmd = &cobra.Command{
	Use:   "start [target_url]",
	Short: "Start a reverse proxy to a target URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		port, _ := cmd.Flags().GetInt("port")

		err := proxy.Start(port, target)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	proxyStartCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	proxyCmd.AddCommand(proxyStartCmd)
	rootCmd.AddCommand(proxyCmd)
}
