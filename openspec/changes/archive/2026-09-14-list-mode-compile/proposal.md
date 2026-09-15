## Why

Raw-mode compilation is pure byte concatenation with no structure awareness, so producing valid structured output (e.g. a YAML list of records) depends entirely on every source fragment being hand-indented to match its final position in the concatenated result — a formatting mistake compiles silently into invalid output with no error. Field mode already does the right kind of work (typed resolution, validation, mechanical serialization) but only ever emits a single flat object, so it can't represent a registry of per-project records either. This gap was already anticipated: the archived `target-rename-and-ext` change named "a validating list mode" as a deliberately deferred future possibility when `--ext` was first introduced.

## What Changes

- New declared target mode `list`, alongside today's `raw` and `field`, selected via `bldoc new <target> --mode list`.
- `add-dep`/`rm-dep` on a list-mode target require the ref form `target:record` or `target:record.field` — reusing the existing single-`:`-split ref grammar (`TargetRef.Field`), just interpreted as `record[.field]` when the target's declared mode is `list`. No new CLI syntax.
- A record-only dep (`target:record`, e.g. `PROJECTS:bldoc`) merges its source's whole top-level document into that record: every top-level key must resolve to a scalar, the same restriction already enforced for a single field-path today, just applied to every key at once instead of one.
- A record.field dep (`target:record.field`, e.g. `PROJECTS:bldoc.description`) resolves exactly like today's field-mode dependency — whole-file content, or a dotted field-path with an optional `--format` template — scoped into that one record's field.
- Compiled output for a list-mode target is an array, one entry per record in first-appearance order, each entry a flat map of field name to its resolved value (the `--format`-rendered string when given, else the raw resolved value). Unlike plain field-mode, list-mode fields are NOT wrapped in `{value, raw}` — that wrapper stays exclusive to plain field-mode targets, which keep their existing output shape unchanged.
- `--ext` at `new` is no longer raw-mode-exclusive: also valid with `--mode field` and `--mode list`. For structured (field/list) targets, `--ext` now also selects the serialization encoder, not just the intermediate's filename suffix — `yaml`/`yml` marshals via `gopkg.in/yaml.v3` (already a dependency, used today only for parsing); any other or unset `--ext` keeps today's JSON encoding. Internally, a target's resolved data is held as one format-agnostic representation regardless of source format (already true of field-path parsing) and converted to its output format only at this final write step.
- **BREAKING**: none for existing raw- or field-mode targets that don't opt into `--mode list`. `--ext` with `--mode field` was previously rejected outright at `new`; it is now accepted — a loosening, not a break.

## Capabilities

### New Capabilities

(none — this extends the existing mode system the same way `raw`/`field` already live inside `cli-shell`, `manifest-store`, and `compile-engine`, rather than introducing a separate capability area)

### Modified Capabilities

- `cli-shell`: `new` accepts `--mode list`; `--ext` validity is extended from raw-only to also accept `field`/`list`; `add-dep`/`rm-dep` ref parsing is given a record-aware interpretation of `TargetRef.Field` when the target's declared mode is `list`.
- `manifest-store`: `Target.Mode` gains the valid value `list`; `add-dep`'s mode-exclusivity/addressing checks are extended so a list-mode target requires every dependency to use record-scoped addressing; `Ext` becomes recordable on `field`- and `list`-mode targets, not just `raw`.
- `compile-engine`: new list-mode compilation behavior — merge whole-document record deps, resolve record.field deps, group into one ordered array of records; output serialization for field- and list-mode targets branches on the target's `Ext` (JSON when unset or unrecognized, native YAML via `yaml.v3` when `yaml`/`yml`); a record-only dependency's source is validated the same scalar-only way a single field-path is today, applied across all its top-level keys.

## Impact

- `internal/cli`: `new.go` accepts `--mode list` and loosens the `--ext` restriction to allow `field`/`list`; `ref.go`'s `TargetRef.Field` is interpreted as `record[.field]` downstream when the target is list-mode; `add_dep.go`/`rm_dep.go` validate the ref shape against the target's declared mode.
- `internal/manifest`: `Target.Mode` validation extended to accept `list`; `Ext` validation loosened to allow it alongside `field`/`list` in addition to `raw`.
- `internal/compile`: a new list-mode compilation path (merge, group, order); `Result` gains a shape for list-mode's array-of-records data; output encoding is centralized to choose JSON vs. YAML from the target's `Ext` for both field- and list-mode.
- `internal/cli/make.go`: intermediate-writing logic consults `Ext` for field/list-mode serialization choice, not just for the raw-mode filename suffix.
- `go.mod`: no new dependency — `gopkg.in/yaml.v3` (already present for parsing) is now also used for `Marshal`.
- Not part of this change's own impact, but the motivating consumer: once this ships, `~/.claude`'s `PROJECTS` target and its already-migrated per-project scaffolding (`project-headers/*.yaml`, `entry-description.md` for bldoc, projx, and Repo Rater) get reshaped from hand-formatted raw-mode fragments to plain, ordinarily-indented documents addressed via `record[.field]` deps — a follow-up, not part of this change.
