package cmd

import (
	"fmt"
	"netforge/internal/httpclient"

	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "HTTP testing utilities",
}

var httpGetCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Fetch status and response headers from a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		info, err := httpclient.Get(url)
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
			return
		}
		httpclient.PrintResponseInfo(info)
	},
}

var httpBenchmarkCmd = &cobra.Command{
	Use:   "benchmark [url]",
	Short: "Run a performance benchmark against a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		total, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		fmt.Printf("Benchmarking %s (%d requests, %d concurrent)...\n", url, total, concurrency)
		result := httpclient.Benchmark(url, total, concurrency)
		httpclient.PrintBenchmarkResult(result)
	},
}

func init() {
	httpBenchmarkCmd.Flags().IntP("requests", "n", 100, "Total number of requests to send")
	httpBenchmarkCmd.Flags().IntP("concurrency", "c", 10, "Number of concurrent workers")
	httpCmd.AddCommand(httpGetCmd)
	httpCmd.AddCommand(httpBenchmarkCmd)
	rootCmd.AddCommand(httpCmd)
}
