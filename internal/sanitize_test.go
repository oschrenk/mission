package internal

import (
	"testing"
)

func TestSanitizeHref(t *testing.T) {
	raw := "Link to [Label](http://example.com)"
	sanitized := Sanitize(raw)
	wants := "Link to Label"

	if sanitized != wants {
		t.Fatalf(`%s should equal %s`, sanitized, wants)
	}
}

func TestSanitizeWikiLink(t *testing.T) {
	raw := "Link to [[Wikilink]]"
	sanitized := Sanitize(raw)
	wants := "Link to Wikilink"

	if sanitized != wants {
		t.Fatalf(`%s should equal %s`, sanitized, wants)
	}
}

func TestRelabelWikilink(t *testing.T) {
	raw := "Link to [[Wikilink|Something]]"
	sanitized := Sanitize(raw)
	wants := "Link to Something"

	if sanitized != wants {
		t.Fatalf(`%s should equal %s`, sanitized, wants)
	}
}

func TestSanitizeHashTags(t *testing.T) {
	raw := "Do something #woot"
	sanitized := Sanitize(raw)
	wants := "Do something"

	if sanitized != wants {
		t.Fatalf(`%s should equal %s`, sanitized, wants)
	}
}
