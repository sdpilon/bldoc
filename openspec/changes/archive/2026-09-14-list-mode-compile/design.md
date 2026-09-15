## Context

Grounded in the current implementation:

- `manifest.Target.Mode`/`Ext` are already plain strings recorded per target (`internal/manifest/manifest.go`); `manifest.Dep.Field` is already the raw string captured after a target-ref's `:` (`internal/cli/ref.go`'s `TargetRef.Field`).
- Mode-exclusivity is enforced in exactly one place, `manifest.AddDep` (`internal/manifest/ops.go:65`): it branches on `t.Mode` when declared, otherwise infers from `t.Deps[0].Field != ""`. This is also where any new list-mode addressing rule belongs, for the same reason raw/field exclusivity already lives here rather than in the CLI layer.
- `compile.Target` (`internal/compile/compile.go:32`) dispatches on `target.Deps[0].Field == ""` to choose `compileRaw` vs `compileFields`; `compileFields` resolves each dep via `resolveDepValue` (`internal/compile/compile.go:100`, `internal/compile/source.go`) and always wraps results as `{value, raw}`, always JSON-encoded by `internal/cli/make.go:compileTarget`.
- `parseSource` (`internal/compile/source.go:19`) already normalizes TOML/JSON/YAML into one representation (`map[string]interface{}`) regardless of source format — this is the existing "single internal format, converted only at read time" behavior; this change extends the same principle to the write side.
- See `proposal.md` for why; `specs/*/spec.md` for the full normative behavior. This document covers only how.

## Goals / Non-Goals

**Goals:**
- Reuse existing mechanisms wherever the shape already fits (ref grammar, `parseSource`'s format-agnostic parsing, the scalar-only restriction) rather than introducing parallel machinery.
- Keep the change additive: no behavior change for any target that doesn't declare `--mode list`, and no manifest schema migration for existing `bldoc.toml` files.

**Non-Goals:**
- Preserving a record's field order from its source document. Go's `map[string]interface{}` (used by `parseSource` today) does not retain key order, and neither `encoding/json` nor `yaml.v3` promise map-key order on `Marshal`. Guaranteeing it would mean replacing `parseSource`'s representation throughout the parse path — out of scope here; the spec deliberately leaves within-record field order unspecified.
- Validating that a *raw*-mode target's `--ext` actually matches its concatenated content's format. That's the narrower idea the archived `target-rename-and-ext` change floated and declined; this change solves the underlying problem a different way (list mode never concatenates untyped bytes in the first place, so there's nothing to mis-format), and leaves raw mode exactly as it is.
- Migrating `~/.claude`'s already-hand-formatted `PROJECTS` target and its per-project scaffolding to list mode. That's real follow-up work once this ships, not part of this change (noted in `proposal.md`'s Impact section).

## Decisions

### 1. Reuse `Dep.Field` to hold `record` or `record.field`; no new manifest struct field
The alternative was adding a `Dep.Record` field alongside `Field`. Rejected: `TargetRef` already captures exactly one string after the ref's single `:` (`parseTargetRef`, `internal/cli/ref.go`), and `Dep.Field` already stores that string verbatim for field-mode. Splitting it on `.` downstream (only when the target's declared `Mode` is `list`) needs zero grammar changes and zero manifest schema changes — `bldoc.toml` written by this change round-trips through existing `Save`/`Load` untouched. A literal `.` in a plain field-mode field name is unaffected: the record/field split only ever runs when `Mode == "list"`.

### 2. `record[.field]` validation lives in `manifest.AddDep`, not the CLI layer
Mirrors where raw/field exclusivity already lives (`internal/manifest/ops.go:65`), for the same reason: it needs the target's declared `Mode`, which requires the loaded manifest — the CLI's `add_dep.go` only ever does ref-grammar-level checks (e.g. today's "`--format` requires `:field`") before loading the manifest. Both new checks — "list-mode dep must be record-scoped" and "`--format` requires `record.field`, not bare `record`" — extend `AddDep`'s existing branch on `t.Mode`, adding a `case "list"` alongside the existing raw/field logic, rather than adding a second validation pass before `manifest.Load`.

### 3. Split validation, don't reuse the field-mode/raw-mode boolean check
`AddDep`'s current check is a single boolean (`newFieldMode := dep.Field != ""`). List mode needs a three-way outcome per dep: reject (empty field, or more than one `.`), record-only (zero dots), or record.field (exactly one dot). This becomes a small helper, e.g. `splitRecordField(field string) (record, field string, ok bool)` in `internal/manifest`, called from `AddDep`'s new `case "list"` branch — kept in `manifest` (not `compile`) since both `AddDep`'s validation and `compile`'s later resolution need the identical split, and `manifest` is the lower-level, dependency-free package of the two (`compile` already imports `manifest`, not the reverse).

### 4. `compile.Result` gains a list-mode shape; field-mode's shape is untouched
Add `IsList bool`, `RecordOrder []string` (record IDs in first-appearance order), and `Records map[string]map[string]interface{}` (record ID → field name → resolved value, unwrapped — no `{value, raw}`) to `Result`, alongside the existing `IsField`/`Fields`. `compileFields`'s existing `{value, raw}` wrapper and behavior stay exactly as they are for plain field-mode; list mode is a new code path (`compileList`), not a variant of `compileFields`, because its output shape is genuinely different (grouped + unwrapped vs. flat + wrapped) and forcing one function to produce both would need a shape-selecting branch at nearly every line.

### 5. Record-only merge reuses `parseSource` plus a new whole-document scalar check
`compileList`'s record-only deps call the existing `parseSource` (already format-dispatching on extension, already producing `map[string]interface{}`) and then validate every top-level value with a scalar check extracted from `resolveFieldPath`'s existing single-value switch (`internal/compile/source.go:63`) into a small shared `isScalar(v interface{}) bool`, applied per-key instead of once. `record.field` deps go through the existing `resolveDepValue`/`resolveFieldPath` path unchanged, called once per field the same way `compileFields` already does.

### 6. Per-record duplicate detection replaces the flat `checkDuplicateFields` for list mode
`checkDuplicateFields` (`internal/compile/compile.go:89`) does one flat pass over `target.Deps` — field mode's shape (one dep, one field) makes that sufficient. List mode's field set isn't known ahead of time for a record-only dep (its fields come from parsing the source), so `compileList` tracks `seen map[string]map[string]bool` (record → field → present) incrementally while resolving each dep in order: a record-only dep's every merged key, and a record.field dep's single field, are each checked-and-marked against that record's `seen` set, erroring on the first collision — whether the collision is two record.field deps, two record-only merges, or one of each. `checkDuplicateFields` itself is untouched; it keeps governing plain field-mode only.

### 7. Encoding (JSON vs. YAML) is a separate step from resolution, not folded into `compile.Target`
Add `compile.Encode(result Result, ext string) ([]byte, error)`: for `IsList`/`IsField` results, marshal via `encoding/json` (default, matching today) or `gopkg.in/yaml.v3` when `ext` is `yaml`/`yml` (case-insensitive); for a plain raw result, return `result.Raw` unchanged (raw mode keeps meaning "`ext` is naming only", per the Non-Goals above). `internal/cli/make.go`'s `compileTarget` calls `compile.Target` then `compile.Encode`, and writes the result — replacing its current inline `json.MarshalIndent` call. This keeps `compile.Target` purely about *resolving* values (already true) and isolates the *serialization* choice to one function, matching the "single internal representation, convert only at the IO boundary" principle this whole change is built on. A `[]map[string]interface{}` (built from `Records`/`RecordOrder` in order) is what actually gets marshaled for list mode; a plain `map[string]interface{}` (rewrapped from `Fields` when needed, or `Fields` itself for JSON) for field mode.

### 8. `--ext`'s CLI-level gating becomes unconditional acceptance, not per-mode branching
`internal/cli/new.go` currently rejects `--ext` outright when `--mode field`. That branch is deleted rather than extended to also allow `list` — with the removal, `--ext` is simply always accepted and stored, and it's `compile.Encode` (Decision 7) that gives it meaning per mode. This matches the `cli-shell`/`manifest-store` delta specs, which replace the old "requires raw" requirement outright rather than growing its exception list.

## Risks / Trade-offs

- **Field order within a record is unspecified** (Non-Goal above) → a list-mode target's YAML output could look "jumbled" rather than mirroring the source document's key order, which matters for a human-edited registry like `PROJECTS.yaml`. Mitigation available later without a spec change: sort a record's field names alphabetically before marshaling, for run-to-run determinism even though it won't match source order — not committing to this now since it's an implementation detail the spec leaves open.
- **Merge-collision is a hard error, not last-write-wins** (Decision 6) → less convenient for an intentional "override a default" pattern (e.g. a shared base file plus a per-record file that overrides one key). Mitigation: the workaround is straightforward (don't declare the same key in both sources) and matches this project's existing philosophy of rejecting field-name ambiguity outright (`checkDuplicateFields` already does this for plain field-mode) rather than silently picking a winner.
- **`--ext` gains dual meaning (filename suffix + encoder selection) for field/list mode, but stays filename-only for raw mode** → a slight inconsistency in what the same flag means depending on mode. Mitigation: this is called out explicitly in both the proposal and this document rather than left implicit; raw mode has no structured representation to encode in the first place, so there's nothing to select there.

## Migration Plan

No manifest schema migration: existing `bldoc.toml` files with `raw`/`field`-mode (or mode-inferred) targets parse and compile identically to today, since `Mode: "list"` and its associated `Dep.Field` interpretation are purely additive. No `.bldoc/` intermediate needs migrating either — `make` always fully recompiles (existing behavior, unchanged).

Rollout is therefore just: ship the code, then (as separate follow-up work, not part of this change) reshape `~/.claude/bldoc.toml`'s `PROJECTS` target from raw-mode-with-hand-formatted-fragments to `--mode list --ext yaml`, re-adding each project's dependencies via `record[.field]` refs, and reshaping `project-headers/*.yaml`/`entry-description.md` from their current raw-concatenation-shaped content to plain, ordinarily-indented documents.
