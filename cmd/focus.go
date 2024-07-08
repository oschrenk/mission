package cmd

import (
	"fmt"

	m "github.com/oschrenk/mission/internal"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(focusCmd)
}

var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Show Focus",
	Run: func(cmd *cobra.Command, args []string) {
		mission := m.DefaultInstance()
		fmt.Println(mission.Focus())
	},
}

func init() {
}
