package internal

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	model "github.com/oschrenk/mission/model"

	"github.com/fsnotify/fsnotify"
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

func (mission *Mission) Watch() {
	fmt.Println("Watching", mission.settings.CalendarDataPath)

	sketchybar := mission.settings.Sketchybar.Path
	args := []string{"--trigger", mission.settings.Sketchybar.Event}
	fmt.Printf("Using `%s %s %s`", sketchybar, args[0], args[1])

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// listen for events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					cmd := exec.Command(sketchybar, args...)
					_, err := cmd.Output()
					if err != nil {
						fmt.Printf("Failed with: %q\n", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(mission.settings.CalendarDataPath)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
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
