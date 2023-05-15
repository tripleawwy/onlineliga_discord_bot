package formatutils

import (
	"github.com/olekukonko/tablewriter"
	"strings"
)

// ResultsToTable converts a slice of results to a table
// From:
// [
//
//	["OL5", "P1", "VfL Spessartschwalben", "1:1", "SG Glückauf Randersacker"],
//	["OL2", "P8", "Seevetaler Jungs", "4:1", "SC Union 06"],
//
// ]
// To:
// +--------+----------+----------------------+--------+------------------------+
// | League | Position |      Home Team       | Result |        Away Team       |
// +--------+----------+----------------------+--------+------------------------+
// |  OL5   |    P1    | VfL Spessartschwalben |  1:1   | SG Glückauf Randersacker |
// |  OL2   |    P8    |   Seevetaler Jungs    |  4:1   |       SC Union 06     |
// +--------+----------+----------------------+--------+------------------------+
func ResultsToTable(results [][]string) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"League", "Position", "Home Team", "Result", "Away Team"})
	table.AppendBulk(results)
	table.Render()

	return tableString.String()
}
