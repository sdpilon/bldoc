## Why

`bootstrap-cli` gave `bldoc` a command surface that parses arguments
correctly but reports every command as not-yet-implemented — there is no
manifest, so `new`, `add-dep`, `rm-dep`, `rm`, `list`, and `show` have
nothing to act on. Before any compile step can exist, the tool needs a
durable, hand-authored record of what targets exist and what each one
depends on. This change introduces that record and wires the existing
commands to read and write it, while explicitly leaving compilation
(`make`) for a later change.

## What Changes

- Introduce the manifest: a single committed file at the project root
  that lists every target and, per target, its dependencies (source
  refs, optional field paths, optional format templates) and its output
  mode (raw or field).
- Implement real behavior for `new`, `add-dep`, `rm-dep`, `rm`, `list`,
  and `show` against the manifest: creating a target, adding/removing a
  dependency, deleting a target, listing all targets, and showing one
  target's recorded dependencies.
- Enforce raw-mode/field-mode exclusivity per target: a target's
  dependencies are either all unnamed (raw mode) or all named via
  `target:field` (field mode); `add-dep` mixing the two on one target is
  a hard error, not a silent merge.
- Define and enforce manifest-level validation: duplicate target names,
  duplicate dependencies on the same target, and a malformed manifest
  file on disk are all reported as errors rather than silently accepted
  or ignored.
- Explicitly out of scope: `make` continues to report not-yet-implemented
  (compiling a target into an intermediate is the next change); the
  intermediate output schema; partial-file/snippet addressing (only
  whole-file deps and `source:path` field-path addressing are in scope);
  any read of `source:path`'s structured content (parsing TOML/JSON/YAML
  to resolve a field path happens at `make` time, not here — this change
  only needs to store the path string).

## Capabilities

### New Capabilities

- `manifest-store`: the manifest file's format, on-disk location,
  schema (targets, per-target deps, addressing, format templates, output
  mode), and the read/write/validation semantics the CLI commands use to
  mutate it.

### Modified Capabilities

- `cli-shell`: `new`, `add-dep`, `rm-dep`, `rm`, `list`, and `show` no
  longer report not-yet-implemented — each now performs its real
  manifest operation. `make`'s not-yet-implemented behavior is
  unchanged and stays that way until the compile-engine change.

## Impact

Adds manifest read/write logic and a defined file format to what is
currently a stub-only CLI. Touches `internal/cli/*` (the six commands
above) and introduces a new `internal/manifest` package (or equivalent);
exact file format (TOML/JSON/YAML) and package layout are this change's
own `design.md` question. No change to `add-dep`'s or `rm-dep`'s
argument grammar (target-ref/source-ref parsing, already implemented in
`bootstrap-cli`, is reused as-is).
