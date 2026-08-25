# README

Track your mission (with [sketchybar](https://github.com/FelixKratz/SketchyBar))

- Displays and counts today's tasks
- Displays macOS focus
- Emits sketchybar event if today's journal entry changes
  - also emits `JOURNAL_ID` as ENV variable, see also [Triggering Custom Events](https://felixkratz.github.io/SketchyBar/config/events#triggering-custom-events) for more details
- Emits sketchybar event if macOS focus changes
- Support json

## Installation

### nix

```bash
nix profile install github:oschrenk/mission
```

Prebuilt binaries come from the `oschrenk` Cachix cache, which the flake offers
as a substituter.

### home-manager

The flake ships a home-manager module that installs mission and generates
`~/.config/mission/config.toml`, so the configuration below is declared in Nix
rather than hand-edited.

```nix
{
  inputs.mission.url = "github:oschrenk/mission";
  # avoids pulling in a second sketchybar build
  inputs.mission.inputs.nixpkgs.follows = "nixpkgs";

  # in your home-manager configuration:
  imports = [ inputs.mission.homeModules.mission ];

  programs.mission = {
    enable = true;

    vault = {
      name = "memex";
      path = "$HOME/Obsidian/memex";
    };

    journals = {
      default.path = "$HOME/Obsidian/memex/40 Journals/Personal";
      work.path    = "$HOME/Obsidian/memex/40 Journals/Work";
    };
  };
}
```

`sketchybar.path` defaults to the `sketchybar` from the same nixpkgs, so it
cannot drift out of date. `sketchybar.eventTask`, `sketchybar.eventFocus` and
`focus.path` default to the values documented below.

### From source

```bash
task install
```

## Configuration

### System

- `mission tasks` might require access to iCloud drive (depending on the location of your Vault or markdown files). macOS **should** prompt you
* `mission focus` and `mission watch` require Full Disk Access so that it can access the user's system file for focus at `$HOME/Library/DoNotDisturb/DB/Assertions.json`. You can do so by going to "System Settings" > "Privacy & Security" > "Full Disk Access". `mission` should be listed if you already executed it once.

If you do use the app with SketchyBar, SketchyBar would need "Full Disk Access" since it would orchestrate the calls.

This needs to be done every time you update `mission`

### Application

`mission` looks for the first configuration file in

- `$XDG_CONFIG_HOME/mission/config.toml`
- `$HOME/.config/mission/config.toml`

Without a configuration file `mission` falls back to the defaults below, which
is enough for `mission focus` but not for `mission tasks`.

You need to point `mission` at your Obsidian vault, which is what task paths
are reported relative to

```
[vault]
name = "memex"
path = "$HOME/Obsidian/memex"
```

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
path = "sketchybar"
event_task = "mission_task"
event_focus = "mission_focus"
```

And where the macOS Focus state is read from (default below)

```
[focus]
path = "$HOME/Library/DoNotDisturb/DB/Assertions.json"
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

