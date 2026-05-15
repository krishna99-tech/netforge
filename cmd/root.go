package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netforge",
	Short: "NetForge - Advanced Networking CLI Toolkit",
	Long: `NetForge is a modular, high-performance network utility tool built with Go, 
featuring port scanning, DNS lookups, HTTP testing, and more.`,
}

var outputFormat string

// Disable mousetrap help text on Windows (which can trigger in some terminals)
func init() {
	cobra.MousetrapHelpText = ""
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json, table)")
}
