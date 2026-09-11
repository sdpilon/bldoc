## Why

There's no way to try `bldoc` end-to-end without pointing it at a real project and risking its files. A small, checked-in example that already exercises both compile modes lets anyone (a new contributor, a future session, the maintainer) run `make` and see real output in seconds, with nothing at stake.

## What Changes

- Add `examples/demo/`, a self-contained playground with a pre-populated `bldoc.toml` and source files, ready to `make` immediately.
- The demo's manifest declares one raw-mode target (`NOTES`, concatenating two whole-file dependencies) and one field-mode target (`SUMMARY`, with fields sourced from TOML, JSON, and YAML — one using a format template, two using raw passthrough) — together exercising every dependency shape and source format `compile-engine` supports.
- Add `examples/demo/README.md` documenting how to run it and the exact `new`/`add-dep` commands that produced the checked-in manifest, so it also demonstrates the CLI construction flow, not just the compiled result.
- Add `examples/demo/.gitignore` for `.bldoc/`, the demo's own compiled output (disposable, always fully recomputed — never worth committing).

## Capabilities

### New Capabilities
- `demo-project`: a checked-in example project under `examples/demo/` that ready-compiles via `bldoc make`, covering raw-mode and field-mode (TOML/JSON/YAML) dependencies, plus a README documenting how it was built and how to explore it.

### Modified Capabilities
(none — this adds example fixtures only; it does not change any CLI, manifest, or compile-engine behavior)

## Impact

- New directory `examples/demo/` (manifest, source files, README, .gitignore). No changes to `cmd/`, `internal/`, or existing specs.
- No new dependencies, no CI wiring — this is a manual playground, not an automated test fixture.
