package internal

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	model "github.com/oschrenk/mission/model"

	"github.com/LNMMusic/optional"
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

func (mission *Mission) Focus() string {
	return GetFocus(mission.settings.Focus.Path)
}

func getJournalFromPath(targetPath string, journals map[string]Journal) optional.Option[Journal] {
	for _, journal := range journals {
		journalDir := filepath.Clean(journal.Path)
		fileDir := filepath.Dir(targetPath)
		if journalDir == fileDir {
			return optional.Some(journal)
		}
	}
	return optional.None[Journal]()
}

func testAndFireTask(path string, journals map[string]Journal, sketchybar Sketchybar) {

	maybeJournal := getJournalFromPath(path, journals)
	if maybeJournal.IsSome() {
		journal := maybeJournal.Unwrap()
		fileName := filepath.Base(path)
		today := time.Now().Local().Format("2006-01-02")
		todayName := fmt.Sprintf("%s.%s", today, journal.Extension)

		if fileName == todayName {

			fmt.Printf("Found change in \"%s\" journal for today: \"%s\"\n", journal.Id, path)

			binary := sketchybar.Path
			eventName := sketchybar.TaskEvent
			journalEnv := fmt.Sprintf("JOURNAL_ID=%s", journal.Id)
			args := [4]string{binary, "--trigger", eventName, journalEnv}

			fmt.Printf("Firing: %s\n", args)
			cmd := exec.Command(binary, args[1:]...)
			_, err := cmd.Output()
			if err != nil {
				fmt.Printf("Failed with: %q\n", err)
			}
		}
	}
}

func (mission *Mission) Watch() {
	journals := mission.settings.Journals
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

				// Obsidian just writes to the target file
				if event.Has(fsnotify.Write) {
					fmt.Printf("Modified file: %s \n", event.Name)

					testAndFireTask(event.Name, journals, sketchybar)

					// DoNotDisturb changes removes, and re-creates the
					//    ~/Library/DoNotDisturb/DB/Assertions.json
				} else if event.Has(fsnotify.Create) {
					// basic test
					if fileName == "Assertions.json" {
						fmt.Printf("Found change in focus file: %s\n", event.Name)
						newFocus := GetFocus(event.Name)

						args := [4]string{sketchybar.Path, "--trigger", sketchybar.FocusEvent, fmt.Sprintf("FOCUS_MODE=%s", newFocus)}
						fmt.Printf("Firing: %s\n", args)
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

	// Add journal paths
	for name, journal := range journals {
		err = watcher.Add(journal.Path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Watching \"%s\" journal at \"%s\"\n", name, journal.Path)
	}

	// Add do not disturb path
	err = watcher.Add(focus.Path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Watching macOS Focus settings at \"%s\"\n", focus.Path)

	// Block main goroutine forever.
	<-make(chan struct{})
}

func getEntryName(granularity Granularity, now time.Time, extension string) string {
	switch granularity {
	case Day:
		return fmt.Sprintf("%s.%s", now.Format("2006-01-02"), extension)
	case Week:
		year, week := now.ISOWeek()
		return fmt.Sprintf("%d-W%02d.%s", year, week, extension)
	case Month:
		return fmt.Sprintf("%s.%s", now.Format("2006-01"), extension)
	case Quarter:
		quarter := (now.Month()-1)/3 + 1
		return fmt.Sprintf("%d-Q%d.%s", now.Year(), quarter, extension)
	case Year:
		return fmt.Sprintf("%d.%s", now.Year(), extension)
	default:
		return fmt.Sprintf("%s.%s", now.Format("2006-01-02"), extension)
	}
}

func (mission *Mission) GetTasksFromJournal(journalName string, granularity Granularity, now time.Time) ([]model.Task, error) {
	journal := mission.settings.Journals[journalName]

	var entry = getEntryName(granularity, now, journal.Extension)
	path := journal.Path + "/" + entry
	return mission.GetTasksFromPath(path)
}

func (mission *Mission) GetTasksFromPath(path string) ([]model.Task, error) {
	data, doc, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	tasks := mission.parseTasks(data, doc)

	return tasks, nil
}
