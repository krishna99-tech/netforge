package cmd

import (
	"fmt"
	tlsaudit "netforge/internal/tls"
	"netforge/internal/utils"
	"strings"

	"github.com/spf13/cobra"
)

var tlsauditCmd = &cobra.Command{
	Use:   "tlsaudit [host]",
	Short: "Full TLS cipher suite enumeration and certificate audit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port, _ := cmd.Flags().GetInt("port")

		fmt.Printf("Auditing TLS configuration for %s:%d...\n", host, port)

		result, err := tlsaudit.Audit(host, port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(result)
		case "table":
			var data [][]string
			data = append(data, []string{"Grade", result.Grade})
			data = append(data, []string{"Certificate", result.Certificate.Subject})
			data = append(data, []string{"Issuer", result.Certificate.Issuer})
			data = append(data, []string{"Days Remaining", fmt.Sprintf("%d", result.Certificate.DaysLeft)})
			for _, v := range result.Versions {
				status := "✗ No"
				if v.Supported {
					status = fmt.Sprintf("✓ Yes  (%s)", v.CipherSuite)
				}
				data = append(data, []string{v.Version, status})
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			grade := result.Grade
			fmt.Printf("\n--- TLS Audit: %s ---\n", result.Host)
			fmt.Printf("Grade:         %s\n", grade)
			fmt.Println(strings.Repeat("-", 50))
			fmt.Printf("Certificate:   %s\n", result.Certificate.Subject)
			fmt.Printf("Issuer:        %s\n", result.Certificate.Issuer)
			fmt.Printf("Expires:       %s (%d days)\n", result.Certificate.ExpiresAt.Format("2006-01-02"), result.Certificate.DaysLeft)
			fmt.Println(strings.Repeat("-", 50))
			fmt.Println("TLS Versions:")
			for _, v := range result.Versions {
				if v.Supported {
					fmt.Printf("  ✓ %-10s Cipher: %s\n", v.Version, v.CipherSuite)
				} else {
					fmt.Printf("  ✗ %-10s Disabled\n", v.Version)
				}
			}
		}
	},
}

func init() {
	tlsauditCmd.Flags().IntP("port", "p", 443, "Target port")
	rootCmd.AddCommand(tlsauditCmd)
}
