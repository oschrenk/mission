# home-manager module for mission.
#
#   inputs.mission.url = "github:oschrenk/mission";
#   imports = [ inputs.mission.homeModules.mission ];
#
#   programs.mission = {
#     enable = true;
#     vault = {
#       name = "memex";
#       path = "$HOME/Obsidian/memex";
#     };
#     journals.default = {
#       path = "$HOME/Obsidian/memex/40 Journals/Personal";
#     };
#   };
#
# Generates $XDG_CONFIG_HOME/mission/config.toml, the first path
# internal.LoadSettings() searches.
#
# `sketchybar.path` points at this flake's pkgs.sketchybar, so set
# `inputs.mission.inputs.nixpkgs.follows = "nixpkgs"` in the consuming flake to
# avoid pulling in a second sketchybar build.
self:
{
  config,
  lib,
  pkgs,
  ...
}:

let
  inherit (lib)
    literalExpression
    mkEnableOption
    mkIf
    mkOption
    types
    ;

  cfg = config.programs.mission;

  tomlFormat = pkgs.formats.toml { };

  journalType = types.submodule {
    options = {
      vault = mkOption {
        type = types.nullOr types.str;
        default = null;
        description = ''
          Vault this journal belongs to. Defaults to {option}`vault.name`.
        '';
      };

      path = mkOption {
        type = types.str;
        example = "$HOME/Obsidian/memex/40 Journals/Personal";
        description = ''
          Directory holding the journal entries. Expanded by mission through
          `os.ExpandEnv`, so `$HOME` may be used literally.
        '';
      };

      extension = mkOption {
        type = types.str;
        default = "md";
        description = "File extension of journal entries.";
      };
    };
  };

  # mission's config type is map[string][]Journal and it only ever reads the
  # first element (internal/settings.go, fromParsed), so each journal is
  # rendered as a single-element array of tables: [[journals.<id>]].
  settings = {
    vault = {
      inherit (cfg.vault) name path;
    };

    journals = lib.mapAttrs (_: journal: [
      {
        vault = if journal.vault == null then cfg.vault.name else journal.vault;
        inherit (journal) path extension;
      }
    ]) cfg.journals;

    sketchybar = {
      inherit (cfg.sketchybar) path;
      event_task = cfg.sketchybar.eventTask;
      event_focus = cfg.sketchybar.eventFocus;
    };

    focus = {
      inherit (cfg.focus) path;
    };
  };
in
{
  options.programs.mission = {
    enable = mkEnableOption "mission, a bridge from Obsidian journals and macOS Focus to sketchybar";

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.mission;
      defaultText = literalExpression "mission.packages.\${system}.mission";
      description = "The mission package to install.";
    };

    vault = {
      name = mkOption {
        type = types.str;
        example = "memex";
        description = "Name of the Obsidian vault, reported as `meta.vault`.";
      };

      path = mkOption {
        type = types.str;
        example = "$HOME/Obsidian/memex";
        description = ''
          Root of the Obsidian vault, used to relativise task paths. Expanded by
          mission through `os.ExpandEnv`, so `$HOME` may be used literally.
        '';
      };
    };

    journals = mkOption {
      type = types.attrsOf journalType;
      default = { };
      example = literalExpression ''
        {
          default.path = "$HOME/Obsidian/memex/40 Journals/Personal";
          work.path    = "$HOME/Obsidian/memex/40 Journals/Work";
        }
      '';
      description = ''
        Journals keyed by the id `mission tasks --journal` selects. One journal
        must be named `default`, which is what `--journal` falls back to.
      '';
    };

    sketchybar = {
      path = mkOption {
        type = types.str;
        default = "${pkgs.sketchybar}/bin/sketchybar";
        defaultText = literalExpression ''"''${pkgs.sketchybar}/bin/sketchybar"'';
        description = ''
          sketchybar binary that events are triggered against. Unlike the other
          paths this one is *not* run through `os.ExpandEnv`, so it has to be
          absolute.
        '';
      };

      eventTask = mkOption {
        type = types.str;
        default = "mission_task";
        description = "Custom sketchybar event fired when a journal's tasks change.";
      };

      eventFocus = mkOption {
        type = types.str;
        default = "mission_focus";
        description = "Custom sketchybar event fired when the macOS Focus mode changes.";
      };
    };

    focus.path = mkOption {
      type = types.str;
      default = "$HOME/Library/DoNotDisturb/DB/Assertions.json";
      description = ''
        macOS Focus state file that `mission focus` reads and `mission watch`
        watches. Expanded by mission through `os.ExpandEnv`.
      '';
    };
  };

  config = mkIf cfg.enable {
    assertions = [
      {
        assertion = cfg.journals == { } || cfg.journals ? default;
        message = "programs.mission: one journal must be named `default`.";
      }
    ];

    home.packages = [ cfg.package ];

    # Written unconditionally: LoadSettings falls back to defaults without a
    # config file, but those defaults have no vault and no journals, so an
    # enabled mission without this file would do nothing useful.
    xdg.configFile."mission/config.toml".source = tomlFormat.generate "mission-config.toml" settings;
  };
}
