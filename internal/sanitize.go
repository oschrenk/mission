package internal

import (
	"regexp"
)

func Sanitize(text string) string {
	// remove links: [text](link) => text
	hrefReg := regexp.MustCompile(`\[(.*?)\][\[\(].*?[\]\)]`)

	// remove wikilinks: [[text]] => text
	wikilinkReg := regexp.MustCompile(`\[\[(.*?)\]\]`)

	text = hrefReg.ReplaceAllString(text, "$1")
	text = wikilinkReg.ReplaceAllString(text, "$1")

	return text
}
