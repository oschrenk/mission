package internal

import (
	"testing"

	model "github.com/oschrenk/mission/model"
)

var defaultSettings = Settings{
	// these paths should not be used during testing
	// CalendarDataPath: ...,
	// Extension:      "md",
}

var mission = NewInstance(defaultSettings)

func TestParseEmptyString(t *testing.T) {
	md := []byte(``)
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 0 {
		t.Fatalf(`len(tasks) should be 0, but was %d`, len(tasks))
	}
}

func TestParseString(t *testing.T) {
	md := []byte(`# Heading`)
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 0 {
		t.Fatalf(`len(tasks) should be 0, but was %d`, len(tasks))
	}
}

func TestParseAsteriskTask(t *testing.T) {
	md := []byte(`* foo`)
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 0 {
		t.Fatalf(`len(tasks) should be 0, but was %d`, len(tasks))
	}
}

func TestParseDashTask(t *testing.T) {
	md := []byte(`- foo`)
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 0 {
		t.Fatalf(`len(tasks) should be 0, but was %d`, len(tasks))
	}
}

func TestParseOpenTask(t *testing.T) {
	md := []byte(`- [ ] foo`)
	want := model.Task{
		State: model.Open,
		Text:  "foo",
		Depth: 0,
	}
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 1 {
		t.Fatalf(`len(tasks) should be 1, but was %d`, len(tasks))
	}

	if tasks[0] != want {
		t.Fatalf(`Wanted %v but got %v`, want, tasks[0])
	}
}

func TestParseClosedTask(t *testing.T) {
	md := []byte(`- [x] foo`)
	want := model.Task{
		State: model.Done,
		Text:  "foo",
		Depth: 0,
	}
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 1 {
		t.Fatalf(`len(tasks) should be 1, but was %d`, len(tasks))
	}

	if tasks[0] != want {
		t.Fatalf(`Wanted %v but got %v`, want, tasks[0])
	}
}

func TestParseCancelledTask(t *testing.T) {
	md := []byte(`- [-] foo`)
	want := model.Task{
		State: model.Cancelled,
		Text:  "foo",
		Depth: 0,
	}
	data, doc, _ := parseString(md)
	tasks := mission.parseTasks(data, doc)

	if len(tasks) != 1 {
		t.Fatalf(`len(tasks) should be 1, but was %d`, len(tasks))
	}

	if tasks[0] != want {
		t.Fatalf(`Wanted %v but got %v`, want, tasks[0])
	}
}
