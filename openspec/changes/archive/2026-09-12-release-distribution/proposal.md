## Why

`bldoc` can currently only be obtained by cloning the repository and
running `go build`/`go install` against source. There is no way to grab
a ready-to-run binary for a given platform, and the binary that does get
built reports a hardcoded `0.0.0-dev` version regardless of what was
actually built. This is friction for anyone who wants to try `bldoc`
without a local Go toolchain, and it makes it impossible to tell which
release a given binary corresponds to.

## What Changes

- Add a GoReleaser configuration that builds `bldoc` binaries for
  linux/macos/windows on amd64/arm64, and packages each as a checksummed
  archive.
- Wire the CLI's version string to the tag GoReleaser builds from
  (via `-ldflags`), instead of the hardcoded `0.0.0-dev` constant.
- Add a GitHub Actions release workflow that runs GoReleaser when a
  `v*` tag is pushed, publishing the built archives and a checksums
  file as GitHub Release assets.
- Document the release process (how to cut a release) and installation
  from a downloaded release archive.

Publishing to package managers (Homebrew, apt, etc.) is out of scope for
this change — GitHub Release assets are the only distribution channel it
adds.

## Capabilities

### New Capabilities
- `release-distribution`: defines how `bldoc` is packaged and published
  as installable release binaries — the GoReleaser build matrix, the
  version string it injects, the tag-triggered release workflow, and the
  release assets it must produce.

### Modified Capabilities

(none — the `cli-shell` capability's `--version` requirement already
covers "prints a version string and exits 0"; this change only affects
*what* string gets built in, not that requirement's behavior.)

## Impact

- New file: `.goreleaser.yaml` (or `.yml`) at the repo root, defining the
  build matrix, archive naming, and checksum generation.
- New file: `.github/workflows/release.yml`, triggered on `v*` tags,
  running GoReleaser with `contents: write` permission to publish to
  GitHub Releases.
- `internal/cli/root.go`: replace the hardcoded `version` constant with a
  package variable set via `-ldflags "-X ...=..."`, defaulting to
  `0.0.0-dev` for local `go build`/`go run` with no ldflags.
- Root-level documentation of the release process and binary installation
  (a root `README.md` is a separate, not-yet-proposed change; this change
  documents release/install steps wherever they land — see `design.md`).
- No changes to `bldoc.toml` manifest format, `make` behavior, or any
  other CLI command surface.
