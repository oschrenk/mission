package internal

import (
	"fmt"
	"time"

	model "github.com/oschrenk/mission/model"
)

type Mission struct {
	settings Settings
}

func DefaultInstance() Mission {
	settings := LoadSettings()
	return NewInstance(settings)
}

func NewInstance(settings Settings) Mission {
	return Mission{settings: settings}
}

func (mission *Mission) GetTasks(dateTime time.Time, tp TimePrecision) ([]model.Task, error) {
	entry := ""
	switch tp {
	case Day:
		entry = fmt.Sprint(dateTime.Format("2006-01-02"), ".", mission.settings.Extension)
	case Week:
		year, week := dateTime.ISOWeek()
		entry = fmt.Sprint(year, "-W", fmt.Sprintf("%02d", week), ".", mission.settings.Extension)
	default:
		return nil, fmt.Errorf("unsupported precision %s", tp)
	}
	path := mission.settings.CalendarDataPath + "/" + entry
	data, doc, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	tasks := mission.parseTasks(data, doc)

	return tasks, nil
}
