## Why

`manifest-io` gave `bldoc` a durable, hand-authored record of every
target and its dependencies, but `make` still reports not-yet-implemented
— there is nothing that turns the manifest's recipe into an actual
intermediate. This change closes that gap: `make` compiles a target's
recorded dependencies into a disposable, mechanically-derived output,
completing the "store → intermediate" pipeline this project exists to
provide.

## What Changes

- Implement `make [<target>]`: compile one target, or every target in
  the manifest if none is given.
- Raw-mode compilation: concatenate a target's whole-file dependencies'
  bytes, in the order they were added, byte-identical to `cat`.
- Field-mode compilation: for each dependency, resolve its value (either
  a whole source file's content, or a field-path resolved from the
  source file parsed as TOML/JSON/YAML), run it through the dependency's
  format template if one was recorded, and emit one JSON object for the
  target with one entry per field, each `{"value": <rendered>, "raw":
  <extracted-pre-format>}`.
- Define the intermediate's on-disk location and naming.
- Always fully recompute from current source state on every `make`
  call — no staleness tracking, no drift detection, safe to re-run any
  time.
- New compile-time validation: reject a target whose recorded
  dependencies use the same target-side field name more than once (the
  JSON schema requires field-name uniqueness; this is not caught by
  `manifest-io`'s existing duplicate check, which only compares
  `(source, path)` pairs).
- Explicitly out of scope: any "verify"/fact-checking step that
  compares the intermediate back against its sources; per-fact
  provenance metadata (already rejected in the architecture vision);
  partial-file/snippet addressing into a source's structured content
  beyond a single scalar field-path.

## Capabilities

### New Capabilities

- `compile-engine`: `make`'s real behavior — raw-mode concatenation,
  field-mode source parsing/field-path resolution/format templating,
  the intermediate's schema and on-disk location, and the associated
  error conditions (unknown target, unsupported source format,
  unresolvable field-path, non-scalar field-path target, duplicate
  field name).

### Modified Capabilities

- `cli-shell`: `make`'s two valid-argument scenarios ("explicit target"
  and "no target") no longer report not-yet-implemented — each now
  compiles for real. The "Valid invocations report not-yet-implemented"
  requirement is removed: after this change, no command exhibits that
  behavior any longer.

## Impact

Adds a new `internal/compile` package (or equivalent) and wires
`internal/cli`'s `make.go` to it. Adds a YAML parsing dependency
alongside the existing TOML one (JSON parsing uses the standard
library). No change to the manifest format itself — `compile-engine`
only reads what `manifest-io` already records.
