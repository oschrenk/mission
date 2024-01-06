package cmd

import (
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

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Show tasks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		m.Logger.Enabled = verbose
		WithSummary, _ := cmd.Flags().GetBool("summary")
		ShowCancelled, _ := cmd.Flags().GetBool("show-cancelled")
		ShowDone, _ := cmd.Flags().GetBool("show-done")

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
		tasks, err := mission.GetTasks(dateTime, precision)
		open := 0
		cancelled := 0
		done := 0
		if err == nil {
			for _, task := range tasks {
				switch task.State {
				case model.Cancelled:
					cancelled = cancelled + 1
					if ShowCancelled {
						fmt.Println(task.String())
					}
				case model.Done:
					done = done + 1
					if ShowDone {
						fmt.Println(task.String())
					}
				case model.Open:
					open = open + 1
					fmt.Println(task.String())
				}
			}
		}

		if !ShowCancelled {
			cancelled = 0
		}

		if WithSummary {
			fmt.Println(summmaryText(open, done, cancelled))
		}

	},
}

func summmaryText(open int, done int, cancelled int) string {
	return fmt.Sprintf("%d/%d tasks", done, open+done+cancelled)
}

func init() {
	tasksCmd.Flags().BoolP("verbose", "v", false, "Log verbose")
	tasksCmd.Flags().BoolP("summary", "s", true, "Print summary")
	tasksCmd.Flags().BoolP("show-cancelled", "c", true, "Show Cancelled")
	tasksCmd.Flags().BoolP("show-done", "d", true, "Show Done")
}
