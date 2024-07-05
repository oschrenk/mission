package internal

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Sketchybar struct {
	Path  string `default:"sketchybar"`
	Event string `default:"mission_watch"`
}

type Settings struct {
	CalendarDataPath string
	Extension        string `default:"md"`
	Sketchybar       Sketchybar
}

func LoadSettings() Settings {

	// set defaults
	viper.SetDefault("journal.extension", "md")
	viper.SetDefault("sketchybar.path", "sketchybar")
	viper.SetDefault("sketchybar.event", "mission_watch")

	// set config type
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// viper traverses paths in order, stops at first
	viper.AddConfigPath("$XDG_CONFIG_HOME/mission")
	viper.AddConfigPath("$HOME/.config/mission")

	// load config
	err := viper.ReadInConfig()
	if err != nil {
		// TODO throw better error
		fmt.Errorf("Error reading config file: %w", err)
	}
	settings := Settings{
		CalendarDataPath: os.ExpandEnv(viper.GetString("journal.path")),
		Extension:        viper.GetString("journal.extension"),
		Sketchybar: Sketchybar{
			Path:  viper.GetString("sketchybar.path"),
			Event: viper.GetString("sketchybar.event"),
		},
	}

	return settings
}
