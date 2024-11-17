package formatutils

import (
	"bytes"
	"github.com/fogleman/gg"
	"github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
	"image"
	"image/png"
	"net/http"
	"os"
	"strings"
)

type ImageConfig struct {
	FontSize    float64
	Margin      float64
	ColWidth    float64
	RowHeight   float64
	BadeWidth   float64
	BadgeHeight float64
	R           float64
	G           float64
	B           float64
}

// MatchResultsToImage converts match results to an image
func MatchResultsToImage(results []parse.MatchResult, config ImageConfig) (bytes.Buffer, error) {
	width, height := calculateImageSize(results, config)
	dc := createNewContext(int(width), int(height))

	if err := loadFontFace(dc, config.FontSize); err != nil {
		return bytes.Buffer{}, err
	}

	dc.SetRGB(config.R, config.G, config.B)
	config.ColWidth = width / 7

	for i, result := range results {
		writeMatchResult(dc, result, config, i)
	}

	buf := new(bytes.Buffer)
	if err := encodeToPNG(dc, buf); err != nil {
		return *buf, err
	}
	return *buf, nil
}

// writeMatchResult writes a single match result to the image
func writeMatchResult(dc *gg.Context, result parse.MatchResult, config ImageConfig, rowIndex int) {
	fields := []string{
		result.LeagueLevel,
		result.BadgeURL,
		result.LeaguePosition,
		result.HomeTeam,
		result.MatchResult,
		result.AwayTeam,
		result.Points,
	}

	for colIndex, field := range fields {
		dc.SetRGB(config.R, config.G, config.B)
		x := config.Margin/2 + config.ColWidth*float64(colIndex)
		y := config.Margin + config.RowHeight*float64(rowIndex)

		if colIndex == 0 {
			// Place the text in the center of the column
			x += config.ColWidth / 2
		}
		if colIndex == 4 {
			setDrawingColor(dc, result.MatchState)
		}

		if isImageURL(field) {
			drawBadge(dc, field, x, y, config)
		} else {
			dc.DrawStringAnchored(field, x, y, 0.5, 0.5)
		}
	}
}

// isImageURL checks if the field is a URL pointing to a PNG image
func isImageURL(field string) bool {
	return strings.HasPrefix(field, "http") && strings.HasSuffix(field, ".png")
}

// drawBadge downloads and draws the badge image
func drawBadge(dc *gg.Context, url string, x, y float64, config ImageConfig) {
	badgeImg, err := downloadImage(url)
	if err != nil {
		return
	}
	x += config.ColWidth * 1 / 4
	y += -13
	badgeImg = resizeImage(badgeImg, int(config.BadeWidth), int(config.BadgeHeight))
	dc.DrawImage(badgeImg, int(x), int(y))
}

// setDrawingColor sets the drawing color based on the match state
func setDrawingColor(dc *gg.Context, matchState string) {
	switch matchState {
	case "WIN":
		dc.SetRGB(0, 1, 0)
	case "LOSS":
		dc.SetRGB(1, 0, 0)
	case "DRAW":
		dc.SetRGB(1, 1, 0)
	default:
		dc.SetRGB(1, 1, 1)
	}
}

func calculateImageSize(results []parse.MatchResult, config ImageConfig) (float64, float64) {
	numRows := len(results)
	if numRows == 0 {
		return 0, 0
	}
	width := calculateImageWidth(results, config)
	height := config.RowHeight*float64(numRows) + config.Margin
	return width, height
}

// calculateImageWidth calculates the width of each column based on the text width of the longest string in each column
func calculateImageWidth(results []parse.MatchResult, config ImageConfig) float64 {
	numColumns := 7 // Number of fields in MatchResult excluding MatchState
	colWidths := make([]float64, numColumns)
	for _, result := range results {
		// Calculate the width of each column based on the text width of the longest string in each column
		colWidths = calculateColumnWidth(result, config)
	}

	// Calculate the total width of the table
	var totalWidth float64
	for _, width := range colWidths {
		totalWidth += width
	}

	return totalWidth
}

func calculateColumnWidth(result parse.MatchResult, config ImageConfig) []float64 {
	fields := []string{
		result.LeagueLevel,
		result.LeaguePosition,
		result.HomeTeam,
		result.MatchResult,
		result.AwayTeam,
		result.Points,
	}

	colWidths := make([]float64, len(fields))
	for i, field := range fields {
		colWidths[i] = calculateTextWidth(field, config) + config.Margin
	}
	return colWidths
}

func calculateTextWidth(text string, config ImageConfig) float64 {
	dc := gg.NewContext(1, 1)
	err := dc.LoadFontFace(os.Getenv("FONT_PATH"), config.FontSize)
	if err != nil {
		return 0
	}
	width, _ := dc.MeasureString(text)
	return width * 2
}

// createNewContext creates a new drawing context with a dark background
func createNewContext(width, height int) *gg.Context {
	dc := gg.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.Clear()
	return dc
}

// loadFontFace loads the font face for the drawing context
func loadFontFace(dc *gg.Context, fontSize float64) error {
	fontPath := os.Getenv("FONT_PATH")
	return dc.LoadFontFace(fontPath, fontSize)
}

// encodeToPNG encodes the image to a PNG format
func encodeToPNG(dc *gg.Context, buf *bytes.Buffer) error {
	return png.Encode(buf, dc.Image())
}

// downloadImage downloads an image from a URL
func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// resizeImage resizes the image to the specified width and height
func resizeImage(img image.Image, width, height int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	ggCtx := gg.NewContextForRGBA(newImg)
	ggCtx.Scale(float64(width)/float64(img.Bounds().Dx()), float64(height)/float64(img.Bounds().Dy()))
	ggCtx.DrawImage(img, 0, 0)

	return newImg
}
