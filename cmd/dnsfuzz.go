package cmd

import (
	"bufio"
	"fmt"
	"netforge/internal/dnsfuzz"
	"netforge/internal/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var dnsfuzzCmd = &cobra.Command{
	Use:   "dnsfuzz",
	Short: "DNS subdomain fuzzer and enumerator",
}

var dnsfuzzBruteCmd = &cobra.Command{
	Use:   "brute [domain]",
	Short: "Brute-force subdomains",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		wordlistFile, _ := cmd.Flags().GetString("wordlist")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		var wordlist []string
		if wordlistFile != "" {
			file, err := os.Open(wordlistFile)
			if err != nil {
				fmt.Printf("Error opening wordlist: %v\n", err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				word := strings.TrimSpace(scanner.Text())
				if word != "" {
					wordlist = append(wordlist, word)
				}
			}
		}

		fmt.Printf("Fuzzing %s with %d concurrent workers...\n", domain, concurrency)
		results, err := dnsfuzz.BruteForce(domain, wordlist, concurrency)
		if err != nil {
			fmt.Printf("Error during fuzzing: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("No subdomains found.")
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(results)
		case "table":
			var data [][]string
			for _, r := range results {
				data = append(data, []string{r.Host, strings.Join(r.IPs, ", ")})
			}
			utils.PrintTable([]string{"Subdomain", "IPs"}, data)
		default:
			fmt.Println("\n--- Discovered Subdomains ---")
			for _, r := range results {
				fmt.Printf("[+] %s -> %v\n", r.Host, r.IPs)
			}
		}
	},
}

func init() {
	dnsfuzzBruteCmd.Flags().StringP("wordlist", "w", "", "Path to custom wordlist file")
	dnsfuzzBruteCmd.Flags().IntP("concurrency", "c", 50, "Number of concurrent workers")
	dnsfuzzCmd.AddCommand(dnsfuzzBruteCmd)
	rootCmd.AddCommand(dnsfuzzCmd)
}
