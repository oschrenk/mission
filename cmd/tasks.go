package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ijt/go-anytime"
	"github.com/spf13/cobra"

	m "github.com/oschrenk/mission/internal"
	model "github.com/oschrenk/mission/model"
)

func init() {
	rootCmd.AddCommand(tasksCmd)
}

type TasksWrapper struct {
	Tasks   []model.Task  `json:"tasks"`
	Summary model.Summary `json:"summary"`
}

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Show tasks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		m.Logger.Enabled = verbose
		WithSummary, _ := cmd.Flags().GetBool("summary")
		ShowCancelled, _ := cmd.Flags().GetBool("show-cancelled")
		AsJson, _ := cmd.Flags().GetBool("json")
		ShowDone, _ := cmd.Flags().GetBool("show-done")
		targetJournal, _ := cmd.Flags().GetString("journal")

		now := time.Now()
		dateTime := now
		precision := m.Day
		if len(args) > 0 {
			rawDate := args[0]
			parsedRange, err := anytime.ParseRange(rawDate, now)
			if err != nil {
				m.Logger.Log(fmt.Sprintf("Failed parsing date \"%s\"", rawDate))
			}
			parsedDay := parsedRange.Time
			precision, err = m.BuildTimePrecision(parsedRange.Duration)
			if err != nil {
				m.Logger.Log(fmt.Sprintf("Can't parse precision \"%s\"", rawDate))
			}
			dateTime = parsedDay
		}

		mission := m.DefaultInstance()
		tasks, err := mission.GetTasks(targetJournal, dateTime, precision)
		open := 0
		cancelled := 0
		done := 0
		filteredTasks := []model.Task{}
		if err == nil {
			for _, task := range tasks {
				switch task.State {
				case model.Cancelled:
					cancelled = cancelled + 1
					if ShowCancelled {
						filteredTasks = append(filteredTasks, task)
					}
				case model.Done:
					done = done + 1
					if ShowDone {
						filteredTasks = append(filteredTasks, task)
					}
				case model.Open:
					open = open + 1
					filteredTasks = append(filteredTasks, task)
				}
			}

			if !ShowCancelled {
				cancelled = 0
			}
		}
		summary := model.Summary{Done: done, Total: open + done + cancelled}
		wrapper := TasksWrapper{filteredTasks, summary}

		if AsJson {
			json, _ := json.MarshalIndent(wrapper, "", "  ")
			fmt.Println(string(json))
		} else {
			for _, task := range filteredTasks {
				fmt.Println(task)
			}

			if WithSummary {
				fmt.Println(summary)
			}
		}

	},
}

func init() {
	tasksCmd.Flags().BoolP("verbose", "v", false, "Log verbose")
	tasksCmd.Flags().BoolP("summary", "s", true, "Print summary")

	tasksCmd.Flags().BoolP("json", "", false, "Print json")
	tasksCmd.Flags().BoolP("show-cancelled", "c", true, "Show Cancelled")
	tasksCmd.Flags().BoolP("show-done", "d", true, "Show Done")
	tasksCmd.Flags().StringP("journal", "j", "default", "Select Journal with id")
}
