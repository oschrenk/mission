package internal

import (
	"regexp"
	"strings"
)

func Sanitize(text string) string {
	// remove links: [text](link) => text
	hrefReg := regexp.MustCompile(`\[(.*?)\][\[\(].*?[\]\)]`)
	text = hrefReg.ReplaceAllString(text, "$1")

	// remove tags: "#foo" => ""
	tagReg := regexp.MustCompile(`\s#.*`)
	text = tagReg.ReplaceAllString(text, "")

	// remove wikilinks: [[text]] => text
	// use label if exists [[text|other]] => other
	wikilinkReg := regexp.MustCompile(`\[\[(.*?)\]\]`)
	matches := wikilinkReg.FindStringSubmatch(text)
	//  found wikilink
	if len(matches) > 0 {
		linkAndMaybeLabel := strings.Split(matches[1], "|")
		if len(linkAndMaybeLabel) == 1 {
			// link found, but no label
			text = wikilinkReg.ReplaceAllString(text, linkAndMaybeLabel[0])
		} else {
			// link with label found, choose first label
			text = wikilinkReg.ReplaceAllString(text, linkAndMaybeLabel[1])
		}
	}

	return text
}
