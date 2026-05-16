package cmd

import (
	"fmt"
	"netforge/internal/snmp"
	"netforge/internal/utils"
	"time"

	"github.com/spf13/cobra"
)

var snmpCmd = &cobra.Command{
	Use:   "snmp",
	Short: "SNMP v1/v2c agent queries",
}

var snmpWalkCmd = &cobra.Command{
	Use:   "walk [host]",
	Short: "Walk an SNMP OID tree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		port, _ := cmd.Flags().GetUint16("port")
		community, _ := cmd.Flags().GetString("community")
		version, _ := cmd.Flags().GetInt("version")
		oid, _ := cmd.Flags().GetString("oid")

		opts := snmp.WalkOptions{
			Target:    target,
			Port:      port,
			Community: community,
			Version:   version,
			Timeout:   5 * time.Second,
			Retries:   3,
			RootOID:   oid,
		}

		fmt.Printf("Walking %s (OID: %s)...\n", target, oid)
		results, err := snmp.Walk(opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(results) == 0 {
			fmt.Println("No results returned.")
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(results)
		case "table":
			var data [][]string
			for _, r := range results {
				data = append(data, []string{r.OID, r.Type, r.Value})
			}
			utils.PrintTable([]string{"OID", "Type", "Value"}, data)
		default:
			fmt.Println("\n--- SNMP Walk Results ---")
			for _, r := range results {
				fmt.Printf("%-40s [%s] = %s\n", r.OID, r.Type, r.Value)
			}
		}
	},
}

var snmpGetCmd = &cobra.Command{
	Use:   "get [host]",
	Short: "Get a specific SNMP OID value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		port, _ := cmd.Flags().GetUint16("port")
		community, _ := cmd.Flags().GetString("community")
		version, _ := cmd.Flags().GetInt("version")
		oid, _ := cmd.Flags().GetString("oid")

		opts := snmp.WalkOptions{
			Target:    target,
			Port:      port,
			Community: community,
			Version:   version,
			Timeout:   5 * time.Second,
			Retries:   3,
		}

		result, err := snmp.Get(opts, oid)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch outputFormat {
		case "json":
			utils.PrintJSON(result)
		default:
			fmt.Printf("OID:   %s\nType:  %s\nValue: %s\n", result.OID, result.Type, result.Value)
		}
	},
}

func init() {
	for _, sub := range []*cobra.Command{snmpWalkCmd, snmpGetCmd} {
		sub.Flags().Uint16P("port", "p", 161, "SNMP agent port")
		sub.Flags().StringP("community", "", "public", "SNMP community string")
		sub.Flags().IntP("version", "", 2, "SNMP version (1 or 2)")
		sub.Flags().StringP("oid", "", "1.3.6.1.2.1", "Root OID to start from")
	}
	snmpCmd.AddCommand(snmpWalkCmd)
	snmpCmd.AddCommand(snmpGetCmd)
	rootCmd.AddCommand(snmpCmd)
}
