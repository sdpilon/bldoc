# release-distribution Specification

## Purpose

The `release-distribution` capability defines how `bldoc` is packaged
and published as installable binaries, so users can download a
ready-to-run binary for their platform instead of building from source.

## Requirements

### Requirement: Release triggered by a version tag

Pushing a tag matching `v*` (e.g. `v0.1.0`) to the repository SHALL
trigger a release build. No other push or pull request SHALL trigger a
release build.

#### Scenario: Tag push triggers a release
- **WHEN** a tag named `v0.1.0` is pushed to the repository
- **THEN** a release build runs and produces release assets for that tag

#### Scenario: Non-tag push does not trigger a release
- **WHEN** a commit is pushed to `main` with no accompanying `v*` tag, or a pull request is opened
- **THEN** no release build runs

### Requirement: Cross-platform build matrix

A release build SHALL produce a `bldoc` binary for each combination of
{linux, macOS, windows} and {amd64, arm64} — six binaries per release.

#### Scenario: All six platform/architecture combinations present
- **WHEN** a release build completes for a given tag
- **THEN** its release assets include a `bldoc` binary for linux/amd64, linux/arm64, macOS/amd64, macOS/arm64, windows/amd64, and windows/arm64

### Requirement: Each platform binary is packaged in a checksummed archive

Each platform/architecture binary SHALL be packaged in its own archive
(`.tar.gz`, or `.zip` for windows), and the release SHALL include a
checksums file covering every archive.

#### Scenario: Archive contains the binary
- **WHEN** a release's linux/amd64 archive is extracted
- **THEN** it contains a `bldoc` binary for that platform/architecture

#### Scenario: Checksums file covers every archive
- **WHEN** a release's checksums file is inspected
- **THEN** it lists one checksum entry for every archive published in that release

### Requirement: Released binaries report the tagged version

A `bldoc` binary built by a release build SHALL report the release's tag
(with any leading `v` stripped) when run with `--version`, rather than a
placeholder or development version string.

#### Scenario: Version matches the release tag
- **WHEN** the `bldoc` binary released for tag `v0.1.0` is run with `--version`
- **THEN** it prints `0.1.0`

#### Scenario: Locally built binary is unaffected
- **WHEN** `bldoc` is built locally with a plain `go build` (no release tooling involved)
- **THEN** it reports a development version string, not a release tag

### Requirement: Release assets attach to a GitHub Release

A release build SHALL publish its archives and checksums file as assets
on a GitHub Release corresponding to the pushed tag.

#### Scenario: Assets attached to the tag's release
- **WHEN** a release build for tag `v0.1.0` completes
- **THEN** a GitHub Release for `v0.1.0` exists with every platform archive and the checksums file attached as downloadable assets
