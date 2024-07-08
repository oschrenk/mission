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
