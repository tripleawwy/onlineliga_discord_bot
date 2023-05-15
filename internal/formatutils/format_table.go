package formatutils

import (
	"bytes"
	"github.com/fogleman/gg"
	"github.com/olekukonko/tablewriter"
	"image"
	"image/png"
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

func DrawTable(table [][]string) (bytes.Buffer, error) {
	// Get the number of rows and columns
	numRows := len(table)
	numCols := len(table[0])

	fontSize := 14.0
	margin := fontSize * 2 // Margin added at the top
	// Calculate the width and height of the table based on the number of rows and columns
	colWidth := 150.0
	rowWidth := fontSize * 2
	width := float64(numCols * int(colWidth))
	height := float64(numRows*int(rowWidth)) + margin

	// Create a new image with dark but not black background
	dc := gg.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, int(width), int(height))))
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.Clear()

	// Set the font properties
	//if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", fontSize); err != nil {
	if err := dc.LoadFontFace("/home/bjoern/.local/share/fonts/Saja Typeworks/TrueType/Cascadia Code/Cascadia_Code_Bold.ttf", fontSize); err != nil {
		panic(err)
	}

	// Set the drawing color to black
	//dc.SetRGB(0, 0, 0)
	// Set the drawing color to white
	dc.SetRGB(1, 1, 1)

	// Write all strings from left to right for each row
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			// When to string is ANSI colored switch the color accordingly
			if strings.Contains(table[i][j], "\u001B[0;32m") {
				// Set the drawing color to green
				dc.SetRGB(0, 1, 0)
				// strip the ANSI color code (start and end)
				table[i][j] = stripANSI(table[i][j])
			} else if strings.Contains(table[i][j], "\u001B[0;31m") {
				// Set the drawing color to red
				dc.SetRGB(1, 0, 0)
				// strip the ANSI color code (start and end)
				table[i][j] = stripANSI(table[i][j])
			} else if strings.Contains(table[i][j], "\u001B[0;33m") {
				// Set the drawing color to yellow
				dc.SetRGB(1, 1, 0)
				// strip the ANSI color code (start and end)
				table[i][j] = stripANSI(table[i][j])
			} else {
				// Set the drawing color to white
				dc.SetRGB(1, 1, 1)
			}
			dc.DrawStringAnchored(table[i][j], float64(j*int(colWidth))+margin, float64(i*int(rowWidth))+margin, 0.5, 0.5)
		}
	}

	// Encode the image to a buffer
	buf := new(bytes.Buffer)
	err := png.Encode(buf, dc.Image())
	if err != nil {
		return *buf, err
	}

	return *buf, nil
}

// strip the ANSI color code (start and end)
func stripANSI(text string) string {
	text = strings.ReplaceAll(text, "\u001B[0;32m", "")
	text = strings.ReplaceAll(text, "\u001B[0m", "")
	text = strings.ReplaceAll(text, "\u001B[0;31m", "")
	text = strings.ReplaceAll(text, "\u001B[0;33m", "")
	return text
}
