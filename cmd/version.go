package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of NetForge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NetForge v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
