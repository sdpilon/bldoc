## Context

`bldoc` is a single Go binary (`cmd/bldoc`) with no CGo dependencies
(confirmed: `go.mod` requires only `cobra`, `pflag`, `mousetrap`,
`BurntSushi/toml`, `yaml.v3` — all pure Go), so cross-compilation needs
no special toolchain. `internal/cli/root.go` currently hardcodes
`const version = "0.0.0-dev"` and passes it to cobra's `Version` field
— there's no existing mechanism for injecting a real version at build
time. CI (`.github/workflows/ci.yml`) already runs on every PR and push
to `main` via `actions/setup-go@v5`, which the release workflow can
mirror for consistent Go version resolution (`go-version-file: go.mod`).

See `proposal.md` — Why / What Changes for motivation and scope.

## Goals / Non-Goals

**Goals:**
- A tagged release (`vX.Y.Z`) produces downloadable, checksummed
  binaries for linux/macOS/windows × amd64/arm64 with no manual steps.
- The released binary's `--version` output matches the tag.
- The release pipeline is a thin, standard GoReleaser setup — no custom
  build scripting to maintain.

**Non-Goals:**
- Package manager distribution (Homebrew, apt, Scoop, etc.) — explicitly
  deferred per the proposal.
- Signing/notarization of binaries (macOS codesigning, Windows
  Authenticode) — not required for a dev CLI tool at this stage.
- Docker image publishing.
- A root-level `README.md` — tracked as a separate, not-yet-proposed
  change. This change's install/release documentation lives in a
  dedicated `RELEASING.md` (repo root) so it doesn't collide with that
  future README's content; the README change can link to it later.
- Automatic version bumping / changelog generation — GoReleaser's
  changelog defaults (commit list since last tag) are enough for now;
  no conventional-commit changelog tooling.

## Decisions

**GoReleaser for the build/package/publish pipeline.**
Alternative considered: hand-rolled `GOOS`/`GOARCH` matrix in the
workflow YAML using `go build` directly, with manual `tar`/`zip` and
`sha256sum` steps. Rejected: GoReleaser is the de facto standard for Go
CLI releases, handles the matrix/archive/checksum/GitHub-Release-upload
steps declaratively in one config file, and needs far less workflow
YAML to maintain. Its default archive naming
(`bldoc_<version>_<os>_<arch>.<ext>`) and checksum file are used as-is
rather than customized.

**Version injection via `-ldflags -X`, defaulting to `0.0.0-dev`.**
`internal/cli/root.go`'s `version` becomes a `var` (not `const`, since
`-X` can only overwrite package-level string vars) initialized to
`"0.0.0-dev"`. GoReleaser's `ldflags` template
(`-X bldoc/internal/cli.version={{.Version}}`) overwrites it at release
build time; a plain local `go build`/`go run`/`go test` leaves the
default untouched, matching the spec's "locally built binary is
unaffected" scenario. GoReleaser strips a leading `v` from the tag
before substituting `{{.Version}}`, so no manual trimming is needed in
code.

**Tag-triggered workflow, separate from `ci.yml`.**
A new `.github/workflows/release.yml` triggers on `push: tags: ['v*']`
only. Alternative considered: extending `ci.yml` with a conditional
release job. Rejected: keeps the fast PR/push feedback loop
(`ci.yml`) fully decoupled from the release path, which only needs to
run rarely and needs elevated (`contents: write`) permissions that
`ci.yml` doesn't otherwise need.

**Archive format: `.tar.gz` for linux/macOS, `.zip` for windows.**
Matches GoReleaser's own default and general platform convention
(Windows users generally expect `.zip`; Windows lacks a built-in `tar`
in older shells).

## Risks / Trade-offs

- **[Risk]** A build failure on one platform/arch combination could
  block the whole release (GoReleaser fails the entire run on any build
  error) → **Mitigation**: none needed beyond normal CI — the build
  matrix is pure-Go with no cgo, so per-platform build failures are
  unlikely; if one occurs it surfaces immediately in the workflow logs
  rather than silently shipping a partial release.
- **[Risk]** `contents: write` permission on the release workflow is a
  broader grant than `ci.yml`'s `contents: read` → **Mitigation**: scope
  the permission to the `release` job only, and scope the workflow's
  trigger to tag pushes only, so it never runs against untrusted PR
  branches.
- **[Trade-off]** No package manager listing means users still have to
  manually download and place the binary on `PATH` → accepted per
  proposal's explicit scope cut; can be revisited as a follow-up change
  once there's real user demand.

## Open Questions

None — install/versioning documentation location, archive formats, and
version-injection mechanism are all decided above.
