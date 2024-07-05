# README

Track your mission

## Configuration

`mission` looks for the first configuration file in

- `$XDG_CONFIG_HOME/mission/config.toml`
- `$HOME/.config/mission/config.toml`

You need to configure the path containing your journal entries

```
[journal]
path = "$HOME/Library/Mobile Documents/iCloud~md~obsidian/Documents/personal/10 Journals/Personal"
extension = "md"
```

You can configure sketchybar (defaults below)

```
[sketchybar]
path = "/opt/homebrew/bin/sketchybar"
event = "mission_watch"
```
## Usage

### `mission tasks`

Print today's tasks

```
mission tasks
󰄴 Unpack luggage
󰝦 Grocery shopping
  󰝦 Cheese
2 open tasks
```

## Installation

**Via Github**

```
git clone git@github.com:oschrenk/mission.git
cd mission

# installs to $GOBIN/mission
task install
```
````
