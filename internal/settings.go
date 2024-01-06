package internal

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Settings struct {
	CalendarDataPath string
	Extension        string `default:"md"`
}

func LoadSettings() Settings {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	// viper travers paths in order, stops at first
	viper.AddConfigPath("$XDG_CONFIG_HOME/mission")
	viper.AddConfigPath("$HOME/.config/mission")
	err := viper.ReadInConfig()
	if err != nil {
		// TODO throw better error
		fmt.Errorf("Error reading config file: %w", err)
	}

	CalendarDataPath := os.ExpandEnv(viper.GetString("journal.path"))
	Extension := viper.GetString("journal.extension")

	settings := Settings{
		CalendarDataPath: CalendarDataPath,
		Extension:        Extension,
	}

	return settings
}
