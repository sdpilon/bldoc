## Why

Raw-mode intermediates always compile to the extensionless `.bldoc/<target>`, which is awkward for a target whose concatenated output is meant to be consumed as a specific format. Correcting a target's name today requires `rm` + `new` + re-adding every dependency, even though a rename is a single-field edit with nothing else to cascade. Separately, a target's raw/field mode is currently *computed* from whatever its first dependency happens to look like — the very first `add-dep` on a fresh target is completely unconstrained, so a mistaken first call silently and permanently decides the target's mode with no acknowledgment. This change also makes mode an explicit, required choice at creation time to close that gap.

## What Changes

- **BREAKING**: `new` now requires a `--mode <raw|field>` flag; `bldoc new <target>` with no `--mode` is rejected. This only affects the CLI invocation — a `bldoc.toml` written before this change, with no stored mode on its targets, continues to load and compile exactly as before (mode is inferred from its deps, unchanged).
- `add-dep` on a target with a declared mode rejects any dependency — including the first — whose field-addressing doesn't match that mode. A target with no declared mode (i.e. one predating this change) keeps today's behavior: unconstrained first dependency, checked against recorded deps from the second one onward.
- `new` accepts an optional `--ext <ext>` flag, valid only together with `--mode raw`; given with `--mode field` it's rejected immediately at `new` time. When set, a raw-mode target's compiled intermediate is written to `.bldoc/<target>.<ext>` instead of `.bldoc/<target>`.
- A declared field-mode target with no dependencies yet compiles to an empty field-mode intermediate (`{}` at `.bldoc/<target>.json`) instead of falling into the raw-mode empty-file case, since its mode is now known even before any dependency exists.
- A new `rename` command: `bldoc rename <old> <new>`, renaming a target in place. Requires `<old>` to exist and `<new>` to not already exist; dependencies, their order, declared mode (if any), and `ext` (if any) are preserved unchanged.
- `rename` does not touch `.bldoc/`: a stale `.bldoc/<old>[.ext]` is left on disk, consistent with `rm`'s existing behavior (which likewise never touches previously-compiled output).
- `show` reports a target's mode (declared if set, otherwise computed from its dependencies, matching today's behavior) and its extension if set.
- Out of scope for this change: no validation that a raw-mode target's concatenated bytes actually parse as the declared `ext`'s format (e.g. valid YAML) — `ext` is naming only, not a new compile mode. A validating "list mode" is a possible future change, deliberately not pursued here.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `manifest-store`: `new` requires an explicit `--mode`; `add-dep`'s raw/field mode-exclusivity check is extended to enforce a declared mode from the first dependency onward, while preserving today's inferred-mode behavior for targets with no declared mode; add the optional `ext` field (raw-mode only); `rename` command semantics; `show` reports mode and extension.
- `compile-engine`: raw-mode output path becomes `.bldoc/<target>.<ext>` when the target has `ext` set; an empty target's compiled output depends on its mode (empty raw bytes for raw/undeclared, `{}` for declared field) rather than always being empty raw bytes.
- `cli-shell`: recognized command surface grows to eight commands (add `rename`); `rename` requires exactly two positional arguments; `new` requires `--mode <raw|field>` and accepts an optional `--ext`, rejected together with `--mode field`.

## Impact

- `internal/cli`: new `rename.go` command file; `new.go` requires `--mode` and gains `--ext` flag parsing, validated against `--mode` before touching the manifest; `root.go` registers the new command.
- `internal/manifest`: `Target` struct gains optional `Mode` and `Ext` fields; new `RenameTarget` operation; `AddDep`'s mode-exclusivity check branches on whether `Mode` is set (declared-mode targets: check from the first dependency; unset: existing inferred-from-deps behavior, unchanged).
- `internal/compile`: raw-mode output path resolution reads `Target.Ext` when non-empty; the empty-target case branches on `Target.Mode` when set, falling back to today's unconditional empty-raw behavior when unset.
- No changes to field-mode compilation of non-empty targets, `list`'s requirements, or any other existing command.
