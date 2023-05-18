package formatutils

import (
	"bytes"
	"github.com/fogleman/gg"
	"github.com/olekukonko/tablewriter"
	"image"
	"image/png"
	"os"
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
func ResultsToImage(table [][]string) (bytes.Buffer, error) {
	// Get the number of rows and columns
	numRows := len(table)
	numCols := len(table[0])

	// Set the font size and margin
	fontSize := 14.0
	margin := fontSize * 2

	// Calculate the width and height of the table based on the number of rows and columns
	colWidth := 150.0
	rowWidth := fontSize * 2
	width := float64(numCols * int(colWidth))
	height := float64(numRows*int(rowWidth)) + margin

	// Create a new image with dark but not black background
	dc := createNewContext(int(width), int(height))

	// Load the font face
	if err := loadFontFace(dc, fontSize); err != nil {
		return bytes.Buffer{}, err
	}

	// Set the drawing color to white
	dc.SetRGB(1, 1, 1)

	// Write all strings from left to right for each row
	writeStrings(dc, table, colWidth, rowWidth, margin)

	// Encode the image to a buffer
	buf := new(bytes.Buffer)
	err := encodeToPNG(dc, buf)
	if err != nil {
		return *buf, err
	}

	return *buf, nil
}

func createNewContext(width, height int) *gg.Context {
	dc := gg.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.Clear()
	return dc
}

func loadFontFace(dc *gg.Context, fontSize float64) error {
	// Font file path
	fontPath := os.Getenv("FONT_PATH")

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return err
	}
	return nil
}

func writeStrings(dc *gg.Context, table [][]string, colWidth, rowWidth, margin float64) {
	numRows := len(table)
	numCols := len(table[0])

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			setDrawingColor(dc, table[i][j])
			text := stripANSI(table[i][j])
			dc.DrawStringAnchored(text, float64(j*int(colWidth))+margin, float64(i*int(rowWidth))+margin, 0.5, 0.5)
		}
	}
}

func setDrawingColor(dc *gg.Context, text string) {
	if strings.Contains(text, "\u001B[0;32m") {
		// Set the drawing color to green
		dc.SetRGB(0, 1, 0)
	} else if strings.Contains(text, "\u001B[0;31m") {
		// Set the drawing color to red
		dc.SetRGB(1, 0, 0)
	} else if strings.Contains(text, "\u001B[0;33m") {
		// Set the drawing color to yellow
		dc.SetRGB(1, 1, 0)
	} else {
		// Set the drawing color to white
		dc.SetRGB(1, 1, 1)
	}
}

func encodeToPNG(dc *gg.Context, buf *bytes.Buffer) error {
	err := png.Encode(buf, dc.Image())
	if err != nil {
		return err
	}
	return nil
}

// strip the ANSI color code (start and end)
func stripANSI(text string) string {
	text = strings.ReplaceAll(text, "\u001B[0;32m", "")
	text = strings.ReplaceAll(text, "\u001B[0m", "")
	text = strings.ReplaceAll(text, "\u001B[0;31m", "")
	text = strings.ReplaceAll(text, "\u001B[0;33m", "")
	return text
}
