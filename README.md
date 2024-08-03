# README

Track your mission (with [sketchybar](https://github.com/FelixKratz/SketchyBar))

- Displays and counts today's tasks
- Displays macOS focus
- Emits sketchybar event if today's journal entry changes
  - also emits `JOURNAL_ID` as ENV variable, see also [Triggering Custom Events](https://felixkratz.github.io/SketchyBar/config/events#triggering-custom-events) for more details
- Emits sketchybar event if macOS focus changes
- Support json

## Configuration

### System

- `mission tasks` might require access to iCloud drive (depending on the location of your Vault or markdown files). macOS **should** prompt you
* `mission focus` and `mission watch` require Full Disk Access so that it can access the user's system file for focus at `$HOME/Library/DoNotDisturb/DB/Assertions.json`. You can do so by going to "System Settings" > "Privacy & Security" > "Full Disk Access". `mission` should be listed if you already executed it once.

If you do use the app with SketchyBar, SketchyBar would need "Full Disk Access" since it would orchestrating the calls.

### Application

`mission` looks for the first configuration file in

- `$XDG_CONFIG_HOME/mission/config.toml`
- `$HOME/.config/mission/config.toml`

You need to configure the path containing your journal entries

One journal MUST be named `default`

```
[[journals.default]]
path = "$HOME/Library/Mobile Documents/iCloud~md~obsidian/Documents/personal"
extension = "md"

[[journals.work]]
path = "$HOME/Library/Mobile Documents/iCloud~md~obsidian/Documents/work/"
extension = "md"
```

You can configure sketchybar (defaults below)

```
[sketchybar]
path = "/opt/homebrew/bin/sketchybar"
event_task = "mission_task"
event_focus = "mission_focus"
```
## Usage

### `mission tasks`

Print today's tasks

```
mission tasks
󰄴 Unpack luggage
󰝦 Grocery shopping
  󰝦 Cheese
1/2 tasks
```

Print today's tasks from "work" journal

```
mission tasks --journal=work
󰄴 Finish ticket 123
󰝦 Do ticket 456
1/2 tasks
```

Print today's tasks from "work" journal as json

```
mission tasks --journal=work --json
{
  "tasks": [
    {
      "state": "done",
      "text": "Finish ticket 123"
    },
    {
      "state": "open",
      "text": "Do ticket 456"
    },
  ],
  "summary": {
    "done": 1
    "total": 2
  }
}
```
```

### `mission fcous`

Return current macOS focus

```
mission focus
com.apple.focus.work
```
Possible return values (for built in focus)

- `com.apple.donotdisturb.mode.default`
- `com.apple.focus.personal-time`
- `com.apple.focus.work`
- `com.apple.sleep.sleep-mode`

### `mission watch`

Watches the default journal for changes in today's notes and for changes in macOS' builtin focus mode

`mission watch` will log file changes to `stdout` and emit sketchybar events

## Installation

**Via Github**

- installs to `$GOBIN/mission`

```
git clone git@github.com:oschrenk/mission.git
cd mission
task install
```

**Via homebrew**

```
brew tap oschrenk/made git@github.com:oschrenk/homebrew-made
brew install oschrenk/made/mission
```

