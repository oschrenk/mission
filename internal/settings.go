package internal

import (
	"fmt"
	"log"
	"os"
	"time"

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

type Vault struct {
	Name string
	Path string
}

type Journal struct {
	Id        string
	Vault     string
	Path      string `mapstructure:"path"`
	Extension string `mapstructure:"extension"`
}

func (journal *Journal) GetEntryName(granularity Granularity, now time.Time) string {
	switch granularity {
	case Day:
		return fmt.Sprintf("%s.%s", now.Format("2006-01-02"), journal.Extension)
	case Week:
		year, week := now.ISOWeek()
		return fmt.Sprintf("%d-W%02d.%s", year, week, journal.Extension)
	case Month:
		return fmt.Sprintf("%s.%s", now.Format("2006-01"), journal.Extension)
	case Quarter:
		quarter := (now.Month()-1)/3 + 1
		return fmt.Sprintf("%d-Q%d.%s", now.Year(), quarter, journal.Extension)
	case Year:
		return fmt.Sprintf("%d.%s", now.Year(), journal.Extension)
	default:
		return fmt.Sprintf("%s.%s", now.Format("2006-01-02"), journal.Extension)
	}
}

func (journal *Journal) GetEntryPath(granularity Granularity, now time.Time) string {
	var entry = journal.GetEntryName(granularity, now)
	return journal.Path + "/" + entry
}

type parsed struct {
	Vault Vault `mapstructure:"vault"`
	// without the []Journal here, we get "expected a map, got 'slice'"
	Journals   map[string][]Journal `mapstructure:"journals"`
	Sketchybar Sketchybar           `mapstructure:"sketchybar"`
	Focus      Focus                `mapstructure:"focus"`
}

type Settings struct {
	Vault      Vault
	Journals   map[string]Journal
	Sketchybar Sketchybar
	Focus      Focus
}

func fromParsed(parsed parsed) Settings {
	var vault = parsed.Vault
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
		Vault:      vault,
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
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var parsed parsed
	err = viper.Unmarshal(&parsed)
	if err != nil {
		// TODO throw better error
		log.Printf("error reading config file: %s", err)
	}
	return fromParsed(parsed)
}
