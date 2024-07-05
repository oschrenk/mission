package cmd

import (
	m "github.com/oschrenk/mission/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch directory",
	Run: func(cmd *cobra.Command, args []string) {
		mission := m.DefaultInstance()
		mission.Watch()
	},
}

func init() {
}
