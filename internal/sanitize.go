package internal

import (
	"regexp"
)

func Sanitize(text string) string {
	// remove links: [text](link) => text
	hrefReg := regexp.MustCompile(`\[(.*?)\][\[\(].*?[\]\)]`)

	// remove wikilinks: [[text]] => text
	wikilinkReg := regexp.MustCompile(`\[\[(.*?)\]\]`)

	// remove tags: "#foo" => ""
	tagReg := regexp.MustCompile(`\s#.*`)

	text = hrefReg.ReplaceAllString(text, "$1")
	text = wikilinkReg.ReplaceAllString(text, "$1")
	text = tagReg.ReplaceAllString(text, "")

	return text
}
