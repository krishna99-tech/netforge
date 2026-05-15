package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// OutputType defines the format of the command output.
type OutputType string

const (
	OutputText  OutputType = "text"
	OutputJSON  OutputType = "json"
	OutputTable OutputType = "table"
)

// PrintJSON prints any data as a pretty-printed JSON.
func PrintJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

// PrintTable prints data as a formatted table.
func PrintTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	
	// Convert header to []any
	anyHeader := make([]any, len(header))
	for i, h := range header {
		anyHeader[i] = h
	}
	table.Header(anyHeader...)

	// Convert data to []any for Bulk
	table.Bulk(data)
	
	table.Render()
}
