## Context

See `proposal.md` for motivation. `bootstrap-cli` already implemented
the command surface: `internal/cli` has one file per command, each
`RunE` currently returning `notImplemented(cmd)` after any argument
parsing; `parseTargetRef`/`parseSourceRef` (target-ref and source-ref
grammar) already exist and are reused unchanged. There is no manifest
package yet.

## Goals / Non-Goals

**Goals:**
- Define a single on-disk manifest format and location.
- Give `new`, `add-dep`, `rm-dep`, `rm`, `list`, and `show` real,
  tested behavior against that manifest.
- Enforce the invariants from the architecture vision: raw/field mode
  exclusivity per target, no duplicate targets or dependencies.

**Non-Goals:**
- Compiling a target into an intermediate (`make` stays stubbed).
- Validating that a `path`/field-path actually resolves inside its
  source file's structured content — that requires parsing the source
  file, which is compile-time work for the next change.
- Any manifest format migration/versioning tooling.
- Concurrency/locking around manifest reads and writes.

## Decisions

**File format and location: TOML at `bldoc.toml` in the working directory.**
TOML has stable, human-readable diffs for arrays of tables (each target
and each dependency is its own array entry), which matters since the
manifest is committed to git. Alternatives considered: JSON (no
comments, noisier diffs for array insertions/removals) and YAML
(indentation-sensitive, easier to corrupt by hand). TOML is a good fit
and is already a familiar format in this project's own examples
(`pyproject.toml` as a source in the `cli-shell` spec).

**Mode (raw vs. field) is computed, not stored.** A target's mode is
derived from whether its recorded dependencies have a field name, not
tracked as a separate explicit flag. Storing an explicit mode alongside
the deps would create two sources of truth that could drift (e.g. a
mode flag left stale after a bug in dependency removal); computing it
from the deps themselves makes the invariant self-enforcing.

**Dependency identity is (source, path).** `add-dep`'s duplicate check
and `rm-dep`'s match both key on the source ref's `(source, path)` pair,
independent of the target-side field name or format string — two
`add-dep` calls with the same source and path but different `--format`
strings are still a duplicate, since it would create ambiguous
provenance for that one target field.

**New `internal/manifest` package.** Owns the TOML schema struct(s),
`Load`/`Save`, and mutation operations (`AddTarget`, `RemoveTarget`,
`AddDep`, `RemoveDep`, plus read accessors for `list`/`show`) — each
validating existence/duplication/mode-exclusivity and returning a typed
error. `internal/cli` command implementations become thin wrappers:
parse args (unchanged), call into `internal/manifest`, format output or
the error to stderr. This keeps manifest logic unit-testable without
going through Cobra, and matches the existing split where `internal/cli`
owns argument grammar only.

**Full read-modify-write per invocation, no locking.** Each command
that mutates state loads the whole manifest, applies one change, and
saves it back. `bldoc` is a single-user interactive CLI; concurrent
invocations racing on the same manifest are out of scope, consistent
with the architecture vision's "no drift detection, keep it simple"
stance.

**Malformed manifest is a hard error, no auto-repair.** A `bldoc.toml`
that fails to parse is reported and the command exits non-zero; nothing
attempts partial recovery or silently discards the unparseable content.

## Risks / Trade-offs

- **TOML array-of-tables order could be disturbed by an out-of-band
  edit** (someone hand-editing `bldoc.toml`, or an unrelated formatting
  tool reordering it) → order matters for `show` output and will matter
  more once raw-mode concatenation exists. Mitigation: this change only
  ever mutates the manifest through the CLI, which preserves order by
  construction; hand-editing isn't a sanctioned workflow. A future
  `bldoc verify`-style check could flag reordering if it becomes a real
  problem in practice.
- **No concurrency safety** → two simultaneous CLI invocations against
  the same manifest could race and corrupt or drop a change.
  Mitigation: accepted for now under the single-user assumption; revisit
  only if bldoc grows a use case with concurrent invocations (e.g. CI
  jobs running in parallel against the same manifest).
- **TOML couples the project to a specific parsing library.** Mitigation:
  isolated entirely inside `internal/manifest`'s `Load`/`Save`; the rest
  of the CLI only ever sees the parsed Go struct, so swapping formats
  later touches one package.
