package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Sketchybar struct {
	Path       string `mapstructure:"path"`
	TaskEvent  string `mapstructure:"event_task"`
	FocusEvent string `mapstructure:"event_focus"`
}

type Focus struct {
	Path string `mapstructure:"path"`
}

type Journal struct {
	Id        string
	Path      string `mapstructure:"path"`
	Extension string `mapstructure:"extension"`
}

type parsed struct {
	// without the []Journal here, we get "expected a map, got 'slice'"
	Journals   map[string][]Journal `mapstructure:"journals"`
	Sketchybar Sketchybar           `mapstructure:"sketchybar"`
	Focus      Focus                `mapstructure:"focus"`
}

type Settings struct {
	Journals   map[string]Journal
	Sketchybar Sketchybar
	Focus      Focus
}

func fromParsed(parsed parsed) Settings {
	var journals = make(map[string]Journal)
	for id, journal := range parsed.Journals {
		j := journal[0]
		j.Path = os.ExpandEnv(j.Path)
		j.Id = id
		journals[id] = j
	}

	sketchybar := parsed.Sketchybar
	focus := Focus{
		Path: os.ExpandEnv(parsed.Focus.Path),
	}

	return Settings{
		Journals:   journals,
		Sketchybar: sketchybar,
		Focus:      focus,
	}
}

func LoadSettings() Settings {

	// set defaults
	viper.SetDefault("sketchybar.path", "/opt/homebrew/bin/sketchybar")
	viper.SetDefault("sketchybar.event_task", "mission_task")
	viper.SetDefault("sketchybar.event_focus", "mission_event")
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
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	var parsed parsed
	err = viper.Unmarshal(&parsed)
	if err != nil {
		// TODO throw better error
		log.Printf("Error reading config file: %s", err)
	}
	return fromParsed(parsed)
}
