package parse

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func Result(doc *goquery.Document) string {
	var result string
	// Search for only the first element with a class of "team-overview-current-match" and return the text
	selection := doc.Find(".team-overview-current-match").First()
	if selection.Length() > 0 {
		// Strip all children from the selection
		selection.Children().Remove()
		// Get the text from the selection and strip all whitespaces
		result = selection.Text()
		result = strings.TrimSpace(result)
	}
	return result
}
