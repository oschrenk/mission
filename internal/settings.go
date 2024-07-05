package internal

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Sketchybar struct {
	Path       string
	TaskEvent  string
	FocusEvent string
}

type Focus struct {
	Path string
}

type Journal struct {
	Path      string
	Extension string
}

type Settings struct {
	Journal    Journal
	Sketchybar Sketchybar
	Focus      Focus
}

func LoadSettings() Settings {

	// set defaults
	viper.SetDefault("journal.extension", "md")
	viper.SetDefault("sketchybar.path", "/opt/homebrew/bin/sketchybar")
	viper.SetDefault("sketchybar.event.task", "mission_task")
	viper.SetDefault("sketchybar.event.focus", "mission_event")
	viper.SetDefault("focus.path", "$HOME/Library/DoNotDisturb/DB/Assertions.json")

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
		log.Fatalf("Error reading config file: %w", err)

	}
	settings := Settings{
		Journal: Journal{
			Path:      os.ExpandEnv(viper.GetString("journal.path")),
			Extension: viper.GetString("journal.extension"),
		},
		Sketchybar: Sketchybar{
			Path:       viper.GetString("sketchybar.path"),
			TaskEvent:  viper.GetString("sketchybar.event.task"),
			FocusEvent: viper.GetString("sketchybar.event.focus"),
		},
		Focus: Focus{
			Path: os.ExpandEnv(viper.GetString("focus.path")),
		},
	}

	return settings
}
