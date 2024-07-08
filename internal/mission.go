package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	model "github.com/oschrenk/mission/model"

	"github.com/fsnotify/fsnotify"
	"github.com/itchyny/gojq"
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

func parseFocusFile(path string) string {
	query, err := gojq.Parse("try .data[].storeAssertionRecords[].assertionDetails.assertionDetailsModeIdentifier")
	if err != nil {
		log.Println(err)
	}
	bytes, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		log.Println(err)
	}
	dynamic := make(map[string]interface{})
	err = json.Unmarshal(bytes, &dynamic)
	if err != nil {
		log.Println(err)
	}
	iter := query.Run(dynamic) // or query.RunWithContext
	var value string
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			log.Println(err)
		}
		value = fmt.Sprintf("%s", v)
	}
	return value
}

func (mission *Mission) Watch() {
	journal := mission.settings.Journals["default"]
	sketchybar := mission.settings.Sketchybar
	focus := mission.settings.Focus

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
				fileName := filepath.Base(event.Name)

				if event.Has(fsnotify.Write) {
					log.Println("Modified file:", event.Name)

					today := time.Now().Local().Format("2006-01-02")
					todayName := fmt.Sprintf("%s.%s", today, journal.Extension)

					if fileName == todayName {
						log.Println("Found change in today's file:", event.Name)
						args := [3]string{sketchybar.Path, "--trigger", sketchybar.TaskEvent}

						log.Println("Firing:", args)
						cmd := exec.Command(args[0], args[1], args[2])
						_, err := cmd.Output()
						if err != nil {
							fmt.Printf("Failed with: %q\n", err)
						}
					}
				}

				// DoNotDisturb changes remove, and re-create the
				//    ~/Library/DoNotDisturb/DB/Assertions.json
				// file
				if event.Has(fsnotify.Create) {
					// basic test
					if fileName == "Assertions.json" {
						log.Println("Found change in focus file:", event.Name)
						focus := parseFocusFile(event.Name)

						args := [4]string{sketchybar.Path, "--trigger", sketchybar.FocusEvent, fmt.Sprintf("FOCUS_MODE=%s", focus)}
						log.Println("Firing:", args)
						cmd := exec.Command(args[0], args[1], args[2], args[3])
						_, err := cmd.Output()
						if err != nil {
							fmt.Printf("Failed with: %q\n", err)
						}
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

	// Add journal path
	err = watcher.Add(journal.Path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching", journal.Path)

	// Add do not disturb path
	err = watcher.Add(focus.Path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching", focus.Path)

	// Block main goroutine forever.
	<-make(chan struct{})
}

func (mission *Mission) GetTasks(dateTime time.Time, tp TimePrecision) ([]model.Task, error) {
	entry := ""
	switch tp {
	case Day:
		entry = fmt.Sprint(dateTime.Format("2006-01-02"), ".", mission.settings.Journals["default"].Extension)
	case Week:
		year, week := dateTime.ISOWeek()
		entry = fmt.Sprint(year, "-W", fmt.Sprintf("%02d", week), ".", mission.settings.Journals["default"].Extension)
	default:
		return nil, fmt.Errorf("unsupported precision %s", tp)
	}
	path := mission.settings.Journals["default"].Path + "/" + entry
	data, doc, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	tasks := mission.parseTasks(data, doc)

	return tasks, nil
}
