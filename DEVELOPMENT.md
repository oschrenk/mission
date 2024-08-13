# Development

**Requirements**

- [go](https://go.dev/) `brew install go`
- [staticcheck]() `go install honnef.co/go/tools/cmd/staticcheck@latest`

**Commands**

- **build** `task build`
- **run** `task run`
- **test** `task test`
- **lint** `task lint`
- `task build` Build project
- `task run` Run example
- `task test` Run tests
- `task lint` Lint
- `task install` Install app in `$GOBIN/`
- `task uninstall` Removed app from `$GOBIN/`
- `task artifacts` Produces artifact in `./`
- `task tag` Pushes git tag from `VERSION`
- `task release` Creates GitHub release from artifacts
- `task sha` Prints hashes from artifacts
- `task clean` Removes build directory `.build`
- `task updates` Find dependency updates

## Release

1. Increase version number in `VERSION`
2. `task release` to tag and push
3. `task sha` to print shas to stdout
4. Make changes in [homebrew-made](https://github.com/oschrenk/homebrew-made) and push
5. `brew update`
6. `brew upgrade`
