{
  description = "Mission - surface Obsidian journal tasks and macOS Focus in sketchybar";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  # Offer prebuilt binaries from the Cachix cache so `nix profile install`
  # downloads instead of compiling. Consumers are prompted to trust these.
  nixConfig = {
    extra-substituters = [ "https://oschrenk.cachix.org" ];
    extra-trusted-public-keys = [
      "oschrenk.cachix.org-1:3JOMfkq2vFiLw4UsCVwzu8kWFBkuS/3DD5AojcO9pks="
    ];
  };

  outputs =
    { self, nixpkgs }:
    let
      # Single source of truth for the version: ./VERSION holds a bare semver
      # (e.g. 0.6.2); the "v" prefix is added by the taskfile release flow.
      version = nixpkgs.lib.fileContents ./VERSION;

      # Darwin only: mission drives sketchybar and watches
      # ~/Library/DoNotDisturb, so a linux build has nothing to talk to.
      # aarch64 only: it is what CI builds, so it is the only system the
      # binary cache is ever populated for.
      systems = [
        "aarch64-darwin"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: rec {
        mission = pkgs.buildGoModule {
          pname = "mission";
          inherit version;
          src = self;

          # Regenerate after changing go.mod/go.sum: set to lib.fakeHash,
          # run `nix build`, then paste the expected hash from the error.
          vendorHash = "sha256-Rlv5zJzEDzuxrgTkDoGjOhuQfH9Mz8S/A+NpNQ9x1p0=";

          # No -X flags: mission has no Version/Commit/BuildDate vars to stamp.
          ldflags = [
            "-s"
            "-w"
          ];

          # `completion` is hidden (cmd/root.go) but not disabled, and it never
          # reaches the config loader, so it is safe to run in the sandbox.
          nativeBuildInputs = [ pkgs.installShellFiles ];
          postInstall = ''
            installShellCompletion --cmd mission \
              --bash <($out/bin/mission completion bash) \
              --zsh <($out/bin/mission completion zsh) \
              --fish <($out/bin/mission completion fish)
          '';

          meta = {
            description = "Surface Obsidian journal tasks and macOS Focus in sketchybar";
            homepage = "https://github.com/oschrenk/mission";
            mainProgram = "mission";
            platforms = nixpkgs.lib.platforms.darwin;
          };
        };
        default = mission;
      });

      apps = forAllSystems (pkgs: rec {
        mission = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.mission}/bin/mission";
        };
        default = mission;
      });

      homeModules = rec {
        mission = import ./nix/home-manager.nix self;
        default = mission;
      };

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go # go, language
            golangci-lint # go, linter runner
            gopls # go, lsp
          ];
        };
      });
    };
}
