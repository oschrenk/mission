package model

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Task struct {
	State TaskState `json:"state"`
	Text  string    `json:"text"`
	// for now we always skip the field in json
	Depth int `json:"-"`
}

func (t Task) String() string {
	var char string
	switch {
	case t.State == Open:
		char = "󰝦" // nf-md-checkbox_blank_circle_outline, \udb81\udf66
	case t.State == Cancelled:
		char = "" // nf-oct-x_circle, \uf52f
	case t.State == Done:
		char = "󰄴" // nf-md-checkbox_marked_circle_outline, \udb80\udd34
	}
	indent := strings.Repeat(" ", t.Depth*2)

	grey := color.New(color.Faint, color.FgWhite).SprintFunc()
	none := color.New().SprintFunc()

	colorize := none
	if t.State == Done {
		colorize = grey
	}

	return fmt.Sprintf("%s%s %s", indent, colorize(char), colorize(t.Text))
}
