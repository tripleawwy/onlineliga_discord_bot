package formatutils

import "strings"

// ColourTextByState returns a string with the given text coloured by the given state
func ColourTextByState(text string, state string) string {
	switch state {
	case "WIN":
		return strings.Join([]string{"\u001b[0;32m", text, "\u001b[0m"}, "")
	case "LOSS":
		return strings.Join([]string{"\u001b[0;31m", text, "\u001b[0m"}, "")
	case "DRAW":
		return strings.Join([]string{"\u001b[0;33m", text, "\u001b[0m"}, "")
	default:
		return text
	}
}
