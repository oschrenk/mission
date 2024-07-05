# Development

**Requirements**

- [air](https://github.com/cosmtrek/air)

**Commands**

- **build** `task build`
- **run** `task run`
- **test** `task test`
- **lint** `task lint`

**Continuous build**

By default `task watch` will continuously watch and rebuild the project upon any change.

To also re-run it you should overwrite the defaults and use

```
air -c ./.air.toml --build.bin "./tmp/main todo"
```

## Release

1. Increase version number in `VERSION`
2. `task release` to tag and push
3. `task sha` to print shas to stdout
4. Make changes in [homebrew-made](https://github.com/oschrenk/homebrew-made) and push
5. `brew update`
6. `brew upgrade`
