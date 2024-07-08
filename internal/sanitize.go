package internal

import (
	"regexp"
)

func Sanitize(text string) string {
	// remove links: [text](link) => text
	var hrefReg = regexp.MustCompile(`\[(.*?)\][\[\(].*?[\]\)]`)

	text = hrefReg.ReplaceAllString(text, "$1")

	return text
}
