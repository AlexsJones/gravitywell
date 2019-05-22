package utils

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func RenderResourceTable(headers []string, data [][]string) {


	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
