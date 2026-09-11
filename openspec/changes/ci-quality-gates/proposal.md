## Why

`bldoc` has three merged changes' worth of Go code and no CI at all —
every build/test/lint check so far has only ever run locally, by
whoever happened to run it. There's nothing stopping a broken build,
failing test, or unformatted file from landing on `main` via a PR.

## What Changes

- Add a GitHub Actions workflow that runs on every pull request and on
  pushes to `main`: `go build`, `go test ./...`, `go vet ./...`, a
  `gofmt` formatting check, and `golangci-lint`.
- Add a `golangci-lint` configuration file with a reasonable default
  linter set for this project.
- No change to `bldoc`'s own CLI behavior, manifest format, or compile
  engine — this is tooling only.

## Capabilities

### New Capabilities

(none — pure tooling, no capability-level behavior change; this change
sets `skip_specs: true` in its `.openspec.yaml`)

### Modified Capabilities

(none)

## Impact

Adds `.github/workflows/ci.yml` and a `.golangci.yml` config. No source
or module changes. Branch protection requiring these checks is a
separate, manual repository-settings step the user will handle once the
checks exist and have run at least once (a required status check has to
have appeared before GitHub will let it be selected).
