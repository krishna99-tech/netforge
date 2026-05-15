package cmd

import (
	"fmt"
	"netforge/internal/monitor"
	"netforge/internal/utils"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "System and network monitoring utilities",
}

var monitorCpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show current CPU usage percentage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Measuring CPU usage (please wait 1 second)...")
		usage, err := monitor.GetCPUUsage()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Current CPU Usage: %.2f%%\n", usage)
	},
}

var monitorPortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "List local listening ports",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching local listening ports...")
		ports, err := monitor.GetListeningPorts()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(ports) == 0 {
			fmt.Println("No listening ports found.")
			return
		}

		fmt.Println("\n--- Local Listening Ports ---")
		fmt.Printf("%-10s %-10s\n", "PORT", "TYPE")
		fmt.Println("----------------------")
		for _, p := range ports {
			fmt.Printf("%-10d %-10s\n", p.Port, p.Type)
		}
	},
}

var monitorSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Show basic system information (OS, Memory, Uptime)",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := monitor.GetSystemInfo()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(info)
		case "table":
			data := [][]string{
				{"OS", info.OS},
				{"Platform", info.Platform},
				{"Kernel", info.KernelVersion},
				{"Total Memory", fmt.Sprintf("%.2f GB", float64(info.TotalMemory)/1024/1024/1024)},
				{"Used Memory", fmt.Sprintf("%.2f GB (%.1f%%)", float64(info.UsedMemory)/1024/1024/1024, float64(info.UsedMemory)/float64(info.TotalMemory)*100)},
				{"Uptime", (time.Duration(info.Uptime) * time.Second).String()},
			}
			utils.PrintTable([]string{"Property", "Value"}, data)
		default:
			fmt.Println("\n--- System Information ---")
			fmt.Printf("OS:             %s\n", info.OS)
			fmt.Printf("Platform:       %s\n", info.Platform)
			fmt.Printf("Kernel:         %s\n", info.KernelVersion)
			fmt.Printf("Total Memory:   %.2f GB\n", float64(info.TotalMemory)/1024/1024/1024)
			fmt.Printf("Used Memory:    %.2f GB (%.1f%%)\n", float64(info.UsedMemory)/1024/1024/1024, float64(info.UsedMemory)/float64(info.TotalMemory)*100)
			fmt.Printf("Uptime:         %v\n", time.Duration(info.Uptime)*time.Second)
		}
	},
}

var monitorBandwidthCmd = &cobra.Command{
	Use:   "bandwidth",
	Short: "Show real-time network bandwidth usage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Measuring bandwidth (please wait 1 second)...")
		info, err := monitor.GetBandwidthUsage()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("\n--- Network Bandwidth ---")
		fmt.Printf("Upload:    %.2f KB/s\n", info.SentRate/1024)
		fmt.Printf("Download:  %.2f KB/s\n", info.RecvRate/1024)
	},
}

var monitorNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Show network interface details and WiFi status",
	Run: func(cmd *cobra.Command, args []string) {
		ssid := monitor.GetWiFiSSID()
		interfaces, err := monitor.GetNetworkInterfaces()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			res := map[string]interface{}{
				"wifi_ssid":  ssid,
				"interfaces": interfaces,
			}
			utils.PrintJSON(res)
		case "table":
			fmt.Printf("Active WiFi SSID: %s\n\n", ssid)
			var data [][]string
			for _, iface := range interfaces {
				ips := strings.Join(iface.IPAddresses, ", ")
				data = append(data, []string{iface.Name, ips, iface.MACAddress})
			}
			utils.PrintTable([]string{"Interface", "IP Addresses", "MAC Address"}, data)
		default:
			fmt.Printf("\nActive WiFi SSID: %s\n", ssid)
			fmt.Println("\n--- Network Interfaces ---")
			for _, iface := range interfaces {
				fmt.Printf("[%s]\n", iface.Name)
				fmt.Printf("  IPs:  %v\n", iface.IPAddresses)
				fmt.Printf("  MAC:  %s\n", iface.MACAddress)
				fmt.Println()
			}
		}
	},
}

func init() {
	monitorCmd.AddCommand(monitorCpuCmd)
	monitorCmd.AddCommand(monitorPortsCmd)
	monitorCmd.AddCommand(monitorSystemCmd)
	monitorCmd.AddCommand(monitorBandwidthCmd)
	monitorCmd.AddCommand(monitorNetworkCmd)
	rootCmd.AddCommand(monitorCmd)
}
