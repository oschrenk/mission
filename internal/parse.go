package internal

import (
	"bytes"
	"regexp"
	"strings"

	model "github.com/oschrenk/mission/model"

	"github.com/LNMMusic/optional"
	"github.com/yuin/goldmark/ast"
)

// - just a list element
// * just a list element
// - [ ] open task
// - [-] cancelled task
// - [x] done task
func (mission *Mission) parseTask(task string, depth int) optional.Option[model.Task] {
	text := strings.TrimSpace(task)
	return mission.parseTaskState(text, depth)
}

func (mission *Mission) parseTaskState(text string, depth int) optional.Option[model.Task] {
	var state model.TaskState

	stateRegEx := `^\[([\s-x])\].*`
	re := regexp.MustCompile(stateRegEx)
	matches := re.FindStringSubmatch(text)

	// if any state trigger is found we assume it's a Todo
	if len(matches) == 2 {
		stateChar := matches[1]
		text = Sanitize(strings.TrimSpace(text[3:]))
		switch {
		case stateChar == model.Open.Trigger():
			state = model.Open
		case stateChar == model.Cancelled.Trigger():
			state = model.Cancelled
		case stateChar == model.Done.Trigger():
			state = model.Done
		}
		// otherwise it's a list item, and we return none
	} else {
		return optional.None[model.Task]()
	}
	task := model.Task{State: state, Text: text, Depth: depth}
	return optional.Some[model.Task](task)
}

func getText(n ast.Node, source []byte) string {
	if n.Type() == ast.TypeBlock {
		var text bytes.Buffer
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			text.Write(line.Value(source))
		}
		return text.String()
	}
	return ""
}

func (mission *Mission) parseTasks(data []byte, doc ast.Node) []model.Task {
	var tasks []model.Task
	var depth = -1
	markerMap := make(map[int]string)

	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if n, ok := node.(*ast.List); ok {
			if enter {
				depth = depth + 1
				markerMap[depth] = string((*n).Marker)
			} else {
				depth = depth - 1
			}
		}

		if n, ok := node.(*ast.ListItem); ok && enter {
			item := n.FirstChild()
			text := getText(item, data)
			maybeTask := mission.parseTask(text, depth)
			if maybeTask.IsSome() {
				tasks = append(tasks, maybeTask.Unwrap())
			}
		}

		return ast.WalkContinue, nil
	})

	return tasks
}
