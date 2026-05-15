package cmd

import (
	"fmt"
	"netforge/internal/whois"

	"github.com/spf13/cobra"
)

var whoisCmd = &cobra.Command{
	Use:   "whois [target]",
	Short: "Fetch WHOIS information for a domain or IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		fmt.Printf("Querying WHOIS for %s...\n", target)

		result, err := whois.Query(target)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(whoisCmd)
}
