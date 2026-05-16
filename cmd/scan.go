package cmd

import (
	"fmt"
	"netforge/internal/scanner"
	"netforge/internal/utils"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [host]",
	Short: "Check for open common ports on a host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
		grabBanner, _ := cmd.Flags().GetBool("banner")
		
		// Common ports
		ports := []int{21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5900, 8000, 8080, 8443 ,3000,5173,8081,8443,8888,2083,2087,2096,2095,4443,8443,8172,7615,7616,3306,3307}

		s := scanner.Scanner{
			Host:           host,
			Ports:          ports,
			MaxConcurrency: maxConcurrency,
			Timeout:        500 * time.Millisecond,
			GrabBanner:     grabBanner,
		}

		foundPorts, err := s.Run()
		if err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			return
		}

		if len(foundPorts) == 0 {
			fmt.Println("No open ports found among the common ones.")
		} else {
			switch outputFormat {
			case "json":
				utils.PrintJSON(foundPorts)
			case "table":
				var data [][]string
				for _, p := range foundPorts {
					banner := p.Banner
					if len(banner) > 30 {
						banner = banner[:27] + "..."
					}
					data = append(data, []string{fmt.Sprintf("%d", p.Port), "OPEN", banner})
				}
				utils.PrintTable([]string{"Port", "State", "Banner"}, data)
			default:
				fmt.Println("\n--- Open Ports ---")
				for _, p := range foundPorts {
					if p.Banner != "" {
						bannerLine := strings.ReplaceAll(p.Banner, "\r\n", " ")
						bannerLine = strings.ReplaceAll(bannerLine, "\n", " ")
						if len(bannerLine) > 50 {
							bannerLine = bannerLine[:47] + "..."
						}
						fmt.Printf("[+] Port %-5d is OPEN  | Banner: %s\n", p.Port, bannerLine)
					} else {
						fmt.Printf("[+] Port %-5d is OPEN\n", p.Port)
					}
				}
			}
		}
	},
}

func init() {
	scanCmd.Flags().IntP("max-concurrency", "c", 100, "Maximum number of concurrent port scans")
	scanCmd.Flags().BoolP("banner", "b", false, "Attempt to grab service banners from open ports")
	rootCmd.AddCommand(scanCmd)
}
