## Why

Recording a target's dependencies one at a time means every list-mode
record, and every multi-field target, takes as many `add-dep` calls as
it has fields — e.g. a `PROJECTS`-style list target with a
`description` and a `path` field per record takes two separate
invocations that must agree on the same record name. `add-dep` should
let several dependencies be recorded in one command.

Designing that batch form surfaced a pre-existing notational
inconsistency: `:` in a ref (`target:field`, `source:path`) has always
meant "address a component within this one ref," never "map this to
that." A batch form needs exactly the second relationship — pairing a
field name with an unrelated source-ref — and reusing `:` for it would
collide with its established meaning. Rather than build the batch form
on a shaky notation, this change also frees `:` for that key-to-value
role by moving ref-internal field addressing to `@`.

## What Changes

- **BREAKING**: the field-addressing symbol inside a target-ref or
  source-ref changes from `:` to `@`. `target:field` becomes
  `target@field`; `source:path` becomes `source@path`; list-mode's
  `target:record.field` becomes `target@record.field`. `source#anchor`
  (and its breadcrumb form) is unchanged — `#` already meant "address a
  location within this source" and keeps that meaning. Every existing
  `add-dep`/`rm-dep` invocation using `:field`/`:path` must be rewritten
  to use `@field`/`@path`.
- `add-dep` gains a second, variadic invocation form: `add-dep <target>
  [<record-name>] <pair> <pair> [<pair>...]`, used when 3 or more
  positional arguments are given. It requires a bare target name (no
  `@`), an explicit record name as the next argument when the target's
  declared mode is `list`, and two or more `pair` arguments after that.
- A `pair` is `<source-ref>` (no colon) for a whole-file, anchor-
  addressed, or list-mode record-only dependency, or
  `<field>:<source-ref>` for a field-mode dependency or a single field
  within a list-mode record. Because a source-ref never itself contains
  `:` (path-addressing now uses `@`), splitting a pair on its first `:`
  is unambiguous in every mode — no per-target-mode parsing rule is
  needed to tell the two pair shapes apart, and the target's existing
  mode-exclusivity checks are what reject a pair whose shape doesn't
  fit the target's mode, exactly as they reject a mismatched dependency
  today.
- The existing two-positional-argument form (`add-dep <target-ref>
  <source-ref>`) is unchanged in shape — it remains the way to record
  exactly one dependency — but its refs now use `@` per the rename
  above.
- `--format` and `--nested`, when given on the batch form, apply
  uniformly to every pair in the invocation. There is no per-pair
  override; a dependency that needs its own format or nested setting is
  recorded with a separate, single-dependency `add-dep` call instead.
- The batch form is atomic: if any pair fails validation (malformed
  ref, unknown target, mode-exclusivity violation, duplicate
  dependency, over-qualified list field, etc.), the CLI reports the
  error and exits non-zero, and the manifest is left completely
  unchanged — no pairs from that invocation are recorded, matching
  every existing `add-dep` failure guarantee.
- The pre-existing ambiguous-anchor warning is still printed per pair,
  same as it is today for a single anchor-addressed dependency.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-shell`: the target-ref/source-ref grammar's field-addressing
  symbol changes from `:` to `@` (breaking), and `add-dep` gains the
  variadic multi-pair invocation form described above, alongside its
  existing two-argument form (now using the renamed symbol).
- `manifest-store`: `add-dep`'s dependency-recording requirements
  (mode exclusivity, duplicate rejection, ambiguous-anchor warning,
  list-mode record/field addressing) now apply per pair within a single
  batch invocation, and that invocation is atomic — a failure on any
  pair leaves the manifest unchanged even though earlier pairs in the
  same call validated successfully. Every scenario referencing
  `:field`/`:path` syntax is updated to `@field`/`@path`.

## Impact

- `internal/cli/ref.go` / `internal/cli/source_ref.go`: the split
  character changes from `:` to `@`; error messages ("at most one ':'
  is allowed") update to name `@`.
- `internal/cli/add_dep.go`: new argument-count branch, record-name
  parsing for list-mode targets, and a `pair` parser (splits on the
  first `:`, then parses the remainder with the renamed `parseSourceRef`).
- `internal/cli/rm_dep.go`: inherits the `@` rename automatically since
  it reuses `parseTargetRef`/`parseSourceRef`; no grammar change of its
  own.
- `internal/manifest/ops.go`: `AddDep` itself is unchanged; the CLI
  loads the manifest once, calls `AddDep` for each pair against that
  same in-memory value, and only calls `Save` once, after every pair has
  succeeded — a failure on any pair skips `Save` entirely, which is
  already sufficient to leave `bldoc.toml` untouched (today's single-
  dependency `add-dep` already works this way: mutate in memory, save
  once at the end).
- No change to `internal/compile/*` — compiled output for a target with
  N dependencies added via one batch call is identical to the same N
  dependencies added via N separate calls, in the same order. Compiled
  output is also unaffected by the `:`/`@` rename, since the manifest
  stores resolved `Field`/`Path` values, not the CLI syntax used to set
  them.
- Every existing `bldoc.toml` on disk is unaffected by the rename (the
  manifest's TOML representation doesn't encode the CLI's `:`/`@`
  choice); only future CLI invocations need the new symbol.
