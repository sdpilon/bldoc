## Why

`--version` currently prints only the release-time semver string, which
stays the fixed placeholder `0.0.0-dev` for any binary built outside the
GoReleaser tag-triggered release flow (e.g. `go install ./cmd/bldoc`
from a branch, to try out a PR before it's merged). That makes it
impossible to tell, from the binary alone, which commit — or whether an
uncommitted change — it was actually built from. Go already embeds this
as VCS build info in every binary built from a version-controlled
checkout (verified via `go version -m`); `--version` should surface it
directly instead of requiring a separate command.

## What Changes

- `--version`'s output gains a parenthetical suffix with the embedded
  VCS commit (short SHA), commit timestamp, and a `modified` flag when
  the working tree had uncommitted changes at build time — read via
  `runtime/debug.ReadBuildInfo()` at startup, the same information
  `go version -m <binary>` already reports.
- This applies uniformly to both dev builds (`go install`/`go build`)
  and real GoReleaser release builds — same code path, no special-casing
  by build method. A real release's `--version` gains the same
  parenthetical alongside its ldflags-injected semver.
- When VCS build info isn't available (e.g. built with `-buildvcs=false`,
  or from a source tree with no `.git`), `--version` SHALL fall back to
  printing exactly what it does today, with no parenthetical and no
  error.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-shell`: "Help and version behave normally" gains the VCS-info
  parenthetical behavior and its availability fallback.

## Impact

- `internal/cli/root.go`: reads `runtime/debug.ReadBuildInfo()` when
  building the root command and appends the parenthetical to the
  reported version string when VCS settings are present.
- No new dependencies (`runtime/debug` is standard library).
- No change to the GoReleaser config or release workflow.
