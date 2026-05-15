package cmd

import (
	"fmt"
	"netforge/internal/ssl"
	"time"

	"github.com/spf13/cobra"
)

var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "SSL certificate utilities",
}

var sslInspectCmd = &cobra.Command{
	Use:   "inspect [host]",
	Short: "Inspect SSL certificate for a host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port, _ := cmd.Flags().GetInt("port")

		fmt.Printf("Inspecting SSL certificate for %s:%d...\n", host, port)
		certs, err := ssl.Inspect(host, port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		for i, cert := range certs {
			fmt.Printf("\n--- Certificate %d ---\n", i+1)
			fmt.Printf("Subject: %s\n", cert.Subject)
			fmt.Printf("Issuer:  %s\n", cert.Issuer)
			fmt.Printf("Expiry:  %v (%s left)\n", cert.Expiry, time.Until(cert.Expiry).Round(time.Hour))
			if len(cert.DNSNames) > 0 {
				fmt.Printf("DNS Names: %v\n", cert.DNSNames)
			}
		}
	},
}

func init() {
	sslInspectCmd.Flags().IntP("port", "p", 443, "Port to connect to")
	sslCmd.AddCommand(sslInspectCmd)
	rootCmd.AddCommand(sslCmd)
}
