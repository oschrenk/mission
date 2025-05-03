package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	m "github.com/oschrenk/mission/internal"
	model "github.com/oschrenk/mission/model"
)

func init() {
	rootCmd.AddCommand(tasksCmd)
}

type Path struct {
	Absolute string `json:"absolute"`
	Relative string `json:"relative"`
}

type Meta struct {
	Vault   string `json:"vault"`
	Journal string `json:"journal"`
	Path    Path   `json:"path"`
}

type TasksWrapper struct {
	Tasks   []model.Task  `json:"tasks"`
	Summary model.Summary `json:"summary"`
	Meta    Meta          `json:"meta"`
}

var granularity m.Granularity = m.Day

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Show tasks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		asJson, _ := cmd.Flags().GetBool("json")
		showCancelled, _ := cmd.Flags().GetBool("show-cancelled")
		showDone, _ := cmd.Flags().GetBool("show-done")
		targetJournal, _ := cmd.Flags().GetString("journal")
		cmd.Flags().Lookup("granularity").Value.String()
		verbose, _ := cmd.Flags().GetBool("verbose")
		withSummary, _ := cmd.Flags().GetBool("summary")

		m.Logger.Enabled = verbose
		mission := m.DefaultInstance()

		var tasks []model.Task
		var err error
		var abs_path string

		if len(args) > 0 {
			if filepath.IsAbs(args[0]) {
				abs_path = args[0]
			} else {
				wd, _ := os.Getwd()
				abs_path = filepath.Join(wd, args[0])
			}
			tasks, err = mission.GetTasksFromPath(abs_path)
		} else {
			journal := mission.Settings.Journals[targetJournal]
			abs_path = journal.GetEntryPath(granularity, time.Now())
			tasks, err = mission.GetTasksFromPath(abs_path)
		}

		open := 0
		cancelled := 0
		done := 0
		filteredTasks := []model.Task{}
		if err == nil {
			for _, task := range tasks {
				switch task.State {
				case model.Cancelled:
					cancelled = cancelled + 1
					if showCancelled {
						filteredTasks = append(filteredTasks, task)
					}
				case model.Done:
					done = done + 1
					if showDone {
						filteredTasks = append(filteredTasks, task)
					}
				case model.Open:
					open = open + 1
					filteredTasks = append(filteredTasks, task)
				}
			}

			if !showCancelled {
				cancelled = 0
			}
		}
		summary := model.Summary{Done: done, Total: open + done + cancelled}
		vault := mission.Settings.Vault
		rel_path, _ := filepath.Rel(vault.Path, abs_path)
		path := Path{abs_path, rel_path}
		fmt.Println(abs_path)
		fmt.Println(vault.Path)
		fmt.Println(rel_path)

		meta := Meta{vault.Name, targetJournal, path}
		wrapper := TasksWrapper{filteredTasks, summary, meta}

		if asJson {
			json, _ := json.MarshalIndent(wrapper, "", "  ")
			fmt.Println(string(json))
		} else {
			for _, task := range filteredTasks {
				fmt.Println(task)
			}

			if withSummary {
				fmt.Println(summary)
			}
		}
	},
}

func init() {
	tasksCmd.Flags().BoolP("json", "", false, "Print json")
	tasksCmd.Flags().Var(
		enumflag.New(&granularity, "granularity", m.GranularityIds, enumflag.EnumCaseInsensitive),
		"granularity",
		"set granularity level; can be 'day', 'week', 'month', 'quarter', 'year'")
	tasksCmd.Flags().BoolP("show-cancelled", "c", true, "Show Cancelled")
	tasksCmd.Flags().BoolP("show-done", "d", true, "Show Done")
	tasksCmd.Flags().BoolP("summary", "s", true, "Print summary")
	tasksCmd.Flags().BoolP("verbose", "v", false, "Log verbose")
	tasksCmd.Flags().StringP("journal", "j", "default", "Select Journal with id")
}
